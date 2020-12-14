package xiaochengxu

const (
	TokenURL             = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	weixin_template_url  = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send"
	mini_subscribe_url   = "https://api.weixin.qq.com/cgi-bin/message/subscribe/send"
	app_pay_unifiedorder = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	refund_url           = "https://api.mch.weixin.qq.com/secapi/pay/refund"
	draw_ulr             = "https://api.mch.weixin.qq.com/mmpaymkttransfers/promotion/transfers"
)

const (
	CacheTokenName = "cacheminitoken"
	TokenTimeOut   = 7100
)
