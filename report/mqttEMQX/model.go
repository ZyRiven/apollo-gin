package mqttemqx

type (
	// 属性上报请求报文
	ReportPropertyReq struct {
		Id      string                 `json:"id"`
		Version string                 `json:"version"`
		Sys     SysInfo                `json:"sys"`
		Params  map[string]interface{} `json:"params"`
		Method  string                 `json:"method"`
	}
	// 标记是否需要回复，1需要回复，0不需要回复
	SysInfo struct {
		Ack int `json:"ack"`
	}
	PropertyNode struct {
		Value      string `json:"value"`
		CreateTime int64  `json:"createTime"`
	}
)
