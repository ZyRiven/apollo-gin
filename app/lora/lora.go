package lora

import (
	"apollo/consts"
	mqttemqx "apollo/report/mqttEMQX"
	"apollo/setting"
	"apollo/utils"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/tarm/serial"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/gpio/gpioreg"
	"periph.io/x/host/v3"
)

var (
	TSerial *serial.Port
	TError  error
	md1Pin  gpio.PinIO
	md0Pin  gpio.PinIO
)

func LoraInit() {
	// 校时 读取 发送
	_, err := host.Init()
	if err != nil {
		log.Fatal(err)
	}
	if _, err := driverreg.Init(); err != nil {
		setting.ZAPS.Errorf("串口初始化失败：%v", err)
	}
	// 定义引脚
	md1Pin = gpioreg.ByName("GPIO73")
	md0Pin = gpioreg.ByName("GPIO75")

	TConfig := &serial.Config{
		Name:        setting.TAddress,
		Baud:        setting.TBaudRate,
		Size:        8,
		Parity:      serial.ParityNone,
		StopBits:    serial.Stop1,
		ReadTimeout: time.Millisecond * time.Duration(setting.TimeOut),
	}

	TSerial, TError = serial.OpenPort(TConfig)
	if TError != nil {
		setting.ZAPS.Errorf("读取串口失败：%v", TError)
	}
	ReadUatr()
	// data := []byte{0xFA ,0x30 ,0x69 ,0x6A ,0x00 ,0x03 ,0x03 ,0x32 ,0x03 ,0x32 ,0x03 ,0xFB}
	// writeUatr(data)

}

func writeUatr(data []byte) {
	if err := md1Pin.Out(gpio.Low); err != nil {
		setting.ZAPS.Errorf("md1Pin拉低失败：%v", err)
	}
	if err := md0Pin.Out(gpio.High); err != nil {
		setting.ZAPS.Errorf("md0Pin拉高失败：%v", err)
	}

	for {
		_, err := TSerial.Write(data)
		if err != nil {
			setting.ZAPS.Debugln("写错误:", err)
		}
		setting.ZAPS.Debugf("发送：%X \n", data)
		time.Sleep(3000 * time.Millisecond)
	}
}

func writeUatrReturn(data []byte) bool {
	if err := md1Pin.Out(gpio.Low); err != nil {
		setting.ZAPS.Errorf("md1Pin拉低失败：%v", err)
	}
	if err := md0Pin.Out(gpio.High); err != nil {
		setting.ZAPS.Errorf("md0Pin拉高失败：%v", err)
	}

	_, err := TSerial.Write(data)
	if err != nil {
		setting.ZAPS.Debugln("写错误:", err)
	}
	setting.ZAPS.Debugln("写入:", data)
	time.Sleep(1000 * time.Millisecond)
	buf := make([]byte, 125)
	n, err := TSerial.Read(buf)
	if err != nil {
		setting.ZAPS.Debugln("读错误:", err)
	}
	setting.ZAPS.Debugf("接收到的报文（十六进制）：% X\n", buf[:n])
	return true
}

// ReadUatr lora读取
func ReadUatr() {
	if err := md1Pin.Out(gpio.Low); err != nil {
		log.Fatal(err)
	}
	if err := md0Pin.Out(gpio.High); err != nil {
		log.Fatal(err)
	}
	go func() {
		for {
			time.Sleep(time.Duration(setting.TReadInterval) * time.Millisecond)
			buf := make([]byte, 125)
			n, err := TSerial.Read(buf)
			if err == nil {
				setting.ZAPS.Debugf("接收到的报文：% X\n", buf[:n])
				massage := fmt.Sprintf("%02X", buf[:n])
				newVar := utils.ValidatorMessage(massage)
				if newVar != "" {
					device := newVar[2:10]
					fmt.Println("设备号：", device)
					t, _ := strconv.ParseInt(newVar[12:14], 16, 32)
					humidity, _ := strconv.ParseInt(newVar[14:16], 16, 32)
					d1, _ := strconv.ParseInt(newVar[16:18], 16, 32)
					d2, _ := strconv.ParseInt(newVar[18:20], 16, 32)
					co := fmt.Sprintf("%.3f", float32(d1*256+d2)/1000000) + "%"
					o, _ := strconv.ParseInt(newVar[20:len(newVar)-2], 16, 32)
					fmt.Println("氧气：", float32(o)/10, " 报文：", newVar[20:len(newVar)-2])
					fmt.Println("温度：", t, " 报文：", newVar[12:14])
					fmt.Println("湿度：", humidity, " 报文：", newVar[14:16])
					fmt.Println("二氧化碳：", co, " 报文：", newVar[16:20])
					data, _ := json.Marshal(&mqttemqx.ReportPropertyReq{
						Id:      uuid.New().String(),
						Version: consts.BuildVersion,
						Sys: mqttemqx.SysInfo{
							Ack: 0,
						},
						Params: map[string]interface{}{
							"t": mqttemqx.PropertyNode{
								Value:      fmt.Sprintf("%v℃", t),
								CreateTime: time.Now().Unix(),
							},
							"h": mqttemqx.PropertyNode{
								Value:      fmt.Sprintf("%v", humidity)+"%RH",
								CreateTime: time.Now().Unix(),
							},
							"co": mqttemqx.PropertyNode{
								Value:      co,
								CreateTime: time.Now().Unix(),
							},
						},
						Method: "thing.event.property.post",
					})
					mqttemqx.MqttPropertyPublish("sensor", fmt.Sprintf("s-%v", device), data)
				}
			}else {
				setting.ZAPS.Debugf("读取失败：%v", err)
			}
		}
	}()
}
