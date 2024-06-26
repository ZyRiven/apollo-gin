package mqttemqx

import (
	"apollo/consts"
	"apollo/setting"
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/google/uuid"
)

var systemMqttClient *MqttWrapperClient

func InitSystemMqtt() error {
	var err error
	systemMqttClient, err = InitMqtt(&MqttConf{
		Addr:     setting.MqttAddr,
		ClientId: setting.MqttClientid,
		// UserName: g.Cfg().MustGet(ctx, "mqtt.auth.userName").String(),
		// Password: g.Cfg().MustGet(ctx, "mqtt.auth.userPassWorld").String(),
	})
	return err
}

func ReportServiceEmqxInit() {
	err := InitSystemMqtt()
	if err != nil {
		setting.ZAPS.Errorf("mqtt连接失败，err：%v", err)
		return
	}
	go ProcessUpLinkFrame()
	ListenGetDevice()
	ListenCommand()
}

func ProcessUpLinkFrame() {
	for {
		data, _ := json.Marshal(&ReportPropertyReq{
			Id:      uuid.New().String(),
			Version: consts.BuildVersion,
			Sys: SysInfo{
				Ack: 0,
			},
			Params: map[string]interface{}{
				"ip": PropertyNode{
					Value:      setting.SystemState.Ip,
					CreateTime: time.Now().Unix(),
				},
				"memUse": PropertyNode{
					Value:      setting.SystemState.MemUse,
					CreateTime: time.Now().Unix(),
				},
				"softVer": PropertyNode{
					Value:      setting.SystemState.SoftVer,
					CreateTime: time.Now().Unix(),
				},
				"runTime": PropertyNode{
					Value:      setting.SystemState.RunTime,
					CreateTime: time.Now().Unix(),
				},
				"diskUse": PropertyNode{
					Value:      setting.SystemState.DiskUse,
					CreateTime: time.Now().Unix(),
				},
				"cpuUse": PropertyNode{
					Value:      setting.SystemState.CpuUse,
					CreateTime: time.Now().Unix(),
				},
			},
			Method: "thing.event.property.post",
		})
		// str := setting.SystemState.Ip + setting.SystemState.MemUse + setting.SystemState.SoftVer +
		// 	setting.SystemState.RunTime + setting.SystemState.DiskUse + setting.SystemState.CpuUse
		// data, _ := json.Marshal(str)
		MqttPropertyPublish(consts.ProductKey, consts.Key, data)
		time.Sleep(25 * time.Second)
	}
}

// MqttPropertyPublish 网关属性上报
func MqttPropertyPublish(productKey, deviceKey string, data []byte) {
	propertyTopic := fmt.Sprintf(consts.PropertyRegisterSubRequestTopic, productKey, deviceKey)
	setting.ZAPS.Infof("上报服务[%s]发布节点上线消息主题%s", "网关属性上报", propertyTopic)
	setting.ZAPS.Debugf("上报服务[%s]发布节点上线消息内容%s", "网关属性上报", data)
	if systemMqttClient.c != nil {
		systemMqttClient.Publish(propertyTopic, data)
	}
}

// 发布订阅主题，获取子设备列表
func ListenGetDevice() {
	propertyTopic := fmt.Sprintf(consts.GetDeviceResponseTopic, consts.ProductKey, consts.Key)
	if err := systemMqttClient.Subscribe(context.Background(), propertyTopic, GetDev); err != nil {
		setting.ZAPS.Errorf("订阅获取设备主题失败：%s", err)
	}
	setting.ZAPS.Infof("EMQX上报服务订阅主题%s成功", consts.GetDeviceResponseTopic)
}

func GetDev(client mqtt.Client, message mqtt.Message) {
	var result map[string]interface{}
	var list []string
	if err := json.Unmarshal(message.Payload(), &result); err != nil {
		setting.ZAPS.Errorf("获取设备列表失败：%s", err)
	}
	for _, k := range result["data"].([]interface{}) {
		sList := strings.Split(k.(string), "-")
		list = append(list, sList[1])
	}
	consts.DeviceList = list
	setting.ZAPS.Debugf("设备列表：%s", consts.DeviceList)
}

// 发布订阅主题，监听发送指令
func ListenCommand() {
	propertyTopic := fmt.Sprintf(consts.SendCommandResponseTopic, consts.ProductKey, consts.Key)
	if err := systemMqttClient.Subscribe(context.Background(), propertyTopic, func(c mqtt.Client, m mqtt.Message) {
		var result map[string]interface{}
		if err := json.Unmarshal(m.Payload(), &result); err != nil {
			setting.ZAPS.Errorf("获取设备列表失败：%s", err)
		}
		strLower := strings.ToLower(result["data"].(string))
		consts.LoraMutex.Lock()
		if strings.HasSuffix(strLower, "fb") && strings.HasPrefix(strLower, "fa") {
			consts.LoraSendList = append(consts.LoraSendList, result["data"].(string))
		} else {
			setting.ZAPS.Errorf("字符串不符合要求: %s", strLower)
			consts.USBSendList = append(consts.USBSendList, result["data"].(string))
		}
		defer consts.LoraMutex.Unlock()
	}); err != nil {
		setting.ZAPS.Errorf("订阅获取设备主题失败：%s", err)
	}
	setting.ZAPS.Infof("EMQX上报服务订阅主题%s成功", consts.GetDeviceResponseTopic)
}
