package mqttemqx

import (
	"apollo/consts"
	"apollo/setting"
	"encoding/json"
	"fmt"
	"time"
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
				"deviceOnline": PropertyNode{
					Value:      setting.SystemState.DeviceOnline,
					CreateTime: time.Now().Unix(),
				},
				"diskUse": PropertyNode{
					Value:      setting.SystemState.DiskUse,
					CreateTime: time.Now().Unix(),
				},
				"devicePacketLoss": PropertyNode{
					Value:      setting.SystemState.DevicePacketLoss,
					CreateTime: time.Now().Unix(),
				},
			},
			Method: "thing.event.property.post",
		})
		MqttPropertyPublish("orangepi-0001", "yjo-0001",data)
		time.Sleep(25 * time.Second)
	}
}

// MqttPropertyPublish 网关属性上报
func MqttPropertyPublish(productKey, deviceKey string,data []byte) {
	propertyTopic := fmt.Sprintf(consts.PropertyRegisterSubRequestTopic, productKey, deviceKey)
	setting.ZAPS.Infof("上报服务[%s]发布节点上线消息主题%s", "网关属性上报", propertyTopic)
	setting.ZAPS.Debugf("上报服务[%s]发布节点上线消息内容%s", "网关属性上报", data)
	if systemMqttClient.c != nil {
		systemMqttClient.Publish(propertyTopic, data)
	}
}
