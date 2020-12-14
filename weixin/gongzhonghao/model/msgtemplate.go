package model

type WeixinTemp struct {
	ToUser     string      `json:"touser"`
	TemplateId string      `json:"template_id"`
	Url        string      `json:"url"`
	Topcolor   string      `json:"topcolor"`
	Data       interface{} `json:"data"`
}

type WeixinTempItem struct {
	Value string `json:"value"`
	Color string `json:"color"`
}
