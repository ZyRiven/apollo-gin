package consts

const (
	//设备上报属性请求topic /sys/${productKey}/${deviceKey}/thing/event/property/post
	PropertyRegisterSubRequestTopic = "/sys/%s/%s/thing/event/property/post"

	//获取子设备 /sys/${productKey}/${deviceKey}/thing/event/property/device/get
	GetDeviceResponseTopic = "/sys/%s/%s/thing/event/property/device/get_reply"
)