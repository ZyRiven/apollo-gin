package report

import (
	"apollo/app/lora"
	mqttemqx "apollo/report/mqttEMQX"
	"apollo/setting"
)

// ReportServiceInit Mqtt初始化
func ReportServiceInit() {
	if setting.MqttOpen {
		mqttemqx.ReportServiceEmqxInit()
		lora.LoraInit()
	}
	//mqttAliyun.ReportServiceAliyunInit()
	//mqttHuawei.ReportServiceHuaweiInit()
	//mqttThingsBoard.ReportServiceThingsBoardInit()
}
