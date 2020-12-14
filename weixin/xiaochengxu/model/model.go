package model

type ErrRsp struct {
	Errmsg  string `json:"errmsg"`
	Errcode int64  `json:"errcode"`
}

type AccessTokenRsp struct {
	Access_token string `json:"access_token"`
	Expires_in   int64  `json:"expires_in"`
	ErrRsp
}

type SubscribeMessage struct {
	Touser            string      `json:"touser"`
	Template_id       string      `json:"template_id"`
	Page              string      `json:"page,omitempty"`
	Miniprogram_state string      `json:"miniprogram_state,omitempty"` //developer为开发版；trial为体验版；formal为正式版；默认为正式版
	Data              interface{} `json:"data"`
}

type MiniToken struct {
	ErrRsp
	Token   string `json:"access_token"`
	Expires int    `json:"expires_in"`
}

type SubcribeValue struct {
	Value string `json:"value"`
}

const SubcribeDateFormat = "2006年01月02日 15:04"