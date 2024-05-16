package consts

import "sync"

const (
	//设备上报属性请求topic /sys/${productKey}/${deviceKey}/thing/event/property/post
	PropertyRegisterSubRequestTopic = "/sys/%s/%s/thing/event/property/post"

	//获取子设备 /sys/${productKey}/${deviceKey}/thing/event/property/device/get
	GetDeviceResponseTopic = "/sys/%s/%s/thing/event/property/device/get_reply"

	//发送指令 /sys/${productKey}/${deviceKey}/thing/event/property/device/get
	SendCommandResponseTopic = "/sys/%s/%s/thing/event/property/send/command_reply"
)

var (
	LoraSendList []string		// lora要发送的数据
	LoraMutex   sync.Mutex		// 加锁
	USBSendList []string		// lora要发送的数据
)
