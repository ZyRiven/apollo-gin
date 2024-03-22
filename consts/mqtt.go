package consts

const (
	//设备上报属性请求topic /sys/${productKey}/${deviceKey}/thing/event/property/post
	PropertyRegisterSubRequestTopic = "/sys/%s/%s/thing/event/property/post"
	//设备上报属性响应topic(平台响应) /sys/${productKey}/${deviceKey}/thing/event/property/post_reply
	PropertyRegisterPubResponseTopic = "/sys/%s/%s/thing/event/property/post_reply"

	//设备上报事件请求topic /sys/${productKey}/${deviceKey}/thing/event/${tsl.event.identifier}/post
	EventRegisterSubRequestTopic = "/sys/%s/%s/thing/event/%s/post"
	//设备上报事件响应topic(平台响应) /sys/${productKey}/${deviceKey}/thing/event/${tsl.event.identifier}_reply
	EventRegisterPubResponseTopic = "/sys/%s/%s/thing/event/%s/post_reply"

	// 设备上报批量属性请求topic /sys/${productKey}/${deviceKey}/thing/event/property/pack/post
	BatchRegisterSubRequestTopic = "/sys/%s/%s/thing/event/property/pack/post"
	// 设备上报批量属性响应topic(平台响应) /sys/${productKey}/${deviceKey}/thing/event/property/pack/post_reply
	BatchRegisterPubResponseTopic = "/sys/%s/%s/thing/event/property/pack/post_reply"

	//设备主动请求配置信息(设备端发起) /sys/${productKey}/${deviceKey}/thing/config/get
	ConfigGetRequestTopic = "/sys/%s/%s/thing/config/get"
	//设备主动请求配置信息(平台响应) /sys/${productKey}/${deviceKey}/thing/config/get_reply
	ConfigGetResponseTopic = "/sys/%s/%s/thing/config/get_reply"
)