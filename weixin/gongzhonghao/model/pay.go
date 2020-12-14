package model

import "encoding/xml"

type RefundReq struct {
	Appid           string `xml:"appid,omitempty"`
	Mchid           string `xml:"mch_id,omitempty"`
	NonceStr        string `xml:"nonce_str,omitempty"`
	Sign            string `xml:"sign,omitempty"`
	Signtype        string `xml:"sign_type,omitempty"`
	Out_trade_no    string `xml:"out_trade_no,omitempty"`
	Transaction_id  string `xml:"transaction_id,omitempty"`
	Out_refund_no   string `xml:"out_refund_no,omitempty"`
	TotalFee        int    `xml:"total_fee,omitempty"`
	RefundFee       int    `xml:"refund_fee,omitempty"`
	Refund_fee_type string `xml:"refund_fee_type,omitempty"`
	Refund_desc     string `xml:"refund_desc,omitempty"`
	Notifyurl       string `xml:"notify_url,omitempty"`
}

type RefundRsp struct {
	XMLName        xml.Name    `xml:"xml,omitempty"`
	Appid          CdataString `xml:"appid,omitempty"`
	Mchid          CdataString `xml:"mch_id,omitempty"`
	NonceStr       CdataString `xml:"nonce_str,omitempty"`
	Sign           CdataString `xml:"sign,omitempty"`
	Resultcode     CdataString `xml:"result_code,omitempty"`
	Errcode        CdataString `xml:"err_code,omitempty"`
	Errcodedes     CdataString `xml:"err_code_des,omitempty"`
	Transaction_id CdataString `xml:"transaction_id,omitempty"`
	Out_trade_no   CdataString `xml:"out_trade_no,omitempty"`
	Out_refund_no  CdataString `xml:"out_refund_no,omitempty"`
	Refund_id      CdataString `xml:"refund_id,omitempty"`
	Refund_fee     int         `xml:"refund_fee,omitempty"`
	Total_fee      int         `xml:"total_fee,omitempty"`
	Cash_fee       int         `xml:"cash_fee,omitempty"`
	ResultCode
}

type AppUnifierOrder struct {
	XMLName        xml.Name `xml:"xml,omitempty"`
	Appid          string   `xml:"appid,omitempty"`
	Mchid          string   `xml:"mch_id,omitempty"`
	Deviceinfo     string   `xml:"device_info,omitempty"`
	NonceStr       string   `xml:"nonce_str,omitempty"`
	Sign           string   `xml:"sign,omitempty"`
	Signtype       string   `xml:"sign_type,omitempty"`
	Body           string   `xml:"body,,omitempty"`
	Detail         string   `xml:"detail,omitempty"`
	Attach         string   `xml:"attach,omitempty"`
	Outtradeno     string   `xml:"out_trade_no,omitempty"`
	Openid         string   `xml:"openid,omitempty"`
	FeeType        string   `xml:"fee_type,omitempty"`
	TotalFee       int      `xml:"total_fee,omitempty"`
	Spbillcreateip string   `xml:"spbill_create_ip,omitempty"`
	Timestart      string   `xml:"time_start,omitempty"`
	Timeexpire     string   `xml:"time_expire,omitempty"`
	Goodstag       string   `xml:"goods_tag,omitempty"`
	Notifyurl      string   `xml:"notify_url,omitempty"`
	Tradetype      string   `xml:"trade_type,omitempty"`
	Limitpay       string   `xml:"limit_pay,omitempty"`
	Sceneinfo      string   `xml:"scene_info,omitempty"`
}

type ResultCode struct {
	Returncode CdataString `xml:"return_code,omitempty"`
	Returnmsg  CdataString `xml:"return_msg,omitempty"`
}

type AppUnifierOrderRsp struct {
	XMLName    xml.Name    `xml:"xml,omitempty"`
	Appid      CdataString `xml:"appid,omitempty"`
	Mchid      CdataString `xml:"mch_id,omitempty"`
	Deviceinfo CdataString `xml:"device_info,omitempty"`
	NonceStr   CdataString `xml:"nonce_str,omitempty"`
	Sign       CdataString `xml:"sign,omitempty"`
	Resultcode CdataString `xml:"result_code,omitempty"`
	Errcode    CdataString `xml:"err_code,omitempty"`
	Errcodedes CdataString `xml:"err_code_des,omitempty"`
	Tradetype  CdataString `xml:"trade_type,omitempty"`
	Prepayid   CdataString `xml:"prepay_id,omitempty"`
	ResultCode
}

//！！！如果微信修改返回字段，验签会失败！！！
//result_code Return_code  both SUCCESS
//代金券ID	coupon_id_$n	否	String(20)	10000	代金券ID,$n为下标，从0开始编号
//单个代金券支付金额	coupon_fee_$n	否	Int	100	单个代金券支付金额,$n为下标，从0开始编号
type NotifyBody struct {
	XMLName              xml.Name `xml:"xml,omitempty"`
	Return_code          string   `xml:"return_code,omitempty"`
	Appid                string   `xml:"appid,omitempty"`
	Bank_type            string   `xml:"bank_type,omitempty"`
	Device_info          string   `xml:"device_info,omitempty"`
	Cash_fee             string   `xml:"cash_fee,omitempty"`
	Fee_type             string   `xml:"fee_type,omitempty"`
	Is_subscribe         string   `xml:"is_subscribe,omitempty"`
	Mch_id               string   `xml:"mch_id,omitempty"`
	Nonce_str            string   `xml:"nonce_str,omitempty"`
	Openid               string   `xml:"openid,omitempty"`
	Out_trade_no         string   `xml:"out_trade_no,omitempty"`
	Result_code          string   `xml:"result_code,omitempty"`
	Sign                 string   `xml:"sign,omitempty"`
	Time_end             string   `xml:"time_end,omitempty"`
	Total_fee            int      `xml:"total_fee,omitempty"`
	Trade_type           string   `xml:"trade_type,omitempty"`
	Transaction_id       string   `xml:"transaction_id,omitempty"`
	Sign_type            string   `xml:"sign_type,omitempty"`
	Settlement_total_fee int      `xml:"settlement_total_fee,omitempty"`
	Err_code             string   `xml:"err_code,omitempty"`
	Err_code_des         string   `xml:"err_code_des,omitempty"`
	Cash_fee_type        string   `xml:"cash_fee_type,omitempty"`
	Coupon_fee           string   `xml:"coupon_fee,omitempty"`
	Coupon_count         string   `xml:"coupon_count,omitempty"`
	Attach               string   `xml:"attach,omitempty"`
	Coupon_id_0          string   `xml:"coupon_id_0,omitempty"`
	Coupon_id_1          string   `xml:"coupon_id_1,omitempty"`
	Coupon_type_0        string   `xml:"coupon_type_0,omitempty"`
	Coupon_type_1        string   `xml:"coupon_type_1,omitempty"`
	Coupon_fee_0         string   `xml:"coupon_fee_0,omitempty"`
	Coupon_fee_1         string   `xml:"coupon_fee_1,omitempty"`
}

type UserDraw struct {
	Mch_appid        string `xml:"mch_appid"`
	Mchid            string `xml:"mchid"`
	Nonce_str        string `xml:"nonce_str"`
	Sign             string `xml:"sign"`
	Partner_trade_no string `xml:"partner_trade_no"`
	Openid           string `xml:"openid"`
	Check_name       string `xml:"check_name"`
	Amount           int    `xml:"amount"`
	Desc             string `xml:"desc"`
	Spbill_create_ip string `xml:"spbill_create_ip"`
}

type UserDrawRsp struct {
	XMLName          xml.Name    `xml:"xml,omitempty"`
	Appid            CdataString `xml:"appid,omitempty"`
	Mchid            CdataString `xml:"mch_id,omitempty"`
	Device_info      string      `xml:"device_info,omitempty"`
	NonceStr         CdataString `xml:"nonce_str,omitempty"`
	Sign             CdataString `xml:"sign,omitempty"`
	Resultcode       CdataString `xml:"result_code,omitempty"`
	Errcode          CdataString `xml:"err_code,omitempty"`
	Errcodedes       CdataString `xml:"err_code_des,omitempty"`
	Partner_trade_no CdataString `xml:"partner_trade_no,omitempty"`
	Payment_no       CdataString `xml:"payment_no,omitempty"`
	Payment_time     CdataString `xml:"payment_time,omitempty"`
	ResultCode
}
type WxConfig struct {
	AppId     string
	Timestamp int64
	NonceStr  string
	Signature string
}

type QrPay struct {
	XMLName        xml.Name `xml:"xml,omitempty"`
	Appid          string   `xml:"appid,omitempty"`
	Mchid          string   `xml:"mch_id,omitempty"`
	Deviceinfo     string   `xml:"device_info,omitempty"`
	NonceStr       string   `xml:"nonce_str,omitempty"`
	Sign           string   `xml:"sign,omitempty"`
	Signtype       string   `xml:"sign_type,omitempty"`
	Body           string   `xml:"body,,omitempty"`
	Detail         string   `xml:"detail,omitempty"`
	Attach         string   `xml:"attach,omitempty"` //商家数据包，原样返回 String(128)
	Outtradeno     string   `xml:"out_trade_no,omitempty"`
	TotalFee       int      `xml:"total_fee,omitempty"`
	FeeType        string   `xml:"fee_type,omitempty"`
	Spbillcreateip string   `xml:"spbill_create_ip,omitempty"`
	GoodsTag       string   `xml:"goods_tag,omitempty"`
	Limitpay       string   `xml:"limit_pay,omitempty"` //no_credit--指定不能使用信用卡支付
	TimeStart      string   `xml:"time_start,omitempty"`
	Timeexpire     string   `xml:"time_expire,omitempty"`
	Receipt        string   `xml:"receipt,omitempty"`
	AuthCode       string   `xml:"auth_code,omitempty"` //付款码
	Sceneinfo      string   `xml:"scene_info,omitempty"`
}

type QrPayRsp struct {
	XMLName              xml.Name `xml:"xml,omitempty"`
	Return_code          string   `xml:"return_code,omitempty"`
	Return_msg           string   `xml:"return_msg,omitempty"`
	Appid                string   `xml:"appid,omitempty"`
	Mch_id               string   `xml:"mch_id,omitempty"`
	Device_info          string   `xml:"device_info,omitempty"`
	Nonce_str            string   `xml:"nonce_str,omitempty"`
	Sign                 string   `xml:"sign,omitempty"`
	Result_code          string   `xml:"result_code,omitempty"`
	Err_code             string   `xml:"err_code,omitempty"`
	Err_code_des         string   `xml:"err_code_des,omitempty"`
	Openid               string   `xml:"openid,omitempty"`
	Is_subscribe         string   `xml:"is_subscribe,omitempty"`
	Trade_type           string   `xml:"trade_type,omitempty"`
	Bank_type            string   `xml:"bank_type,omitempty"`
	Fee_type             string   `xml:"fee_type,omitempty"`
	Total_fee            int      `xml:"total_fee,omitempty"`
	Settlement_total_fee int      `xml:"settlement_total_fee,omitempty"`
	Coupon_fee           string   `xml:"coupon_fee,omitempty"`
	Cash_fee_type        string   `xml:"cash_fee_type,omitempty"`
	Cash_fee             string   `xml:"cash_fee,omitempty"`
	Transaction_id       string   `xml:"transaction_id,omitempty"`
	Out_trade_no         string   `xml:"out_trade_no,omitempty"`
	Attach               string   `xml:"attach,omitempty"`
	Time_end             string   `xml:"time_end,omitempty"`
	Promotion_detail     string   `xml:"promotion_detail,omitempty"`
}

type Orderquery struct {
	XMLName        xml.Name `xml:"xml,omitempty"`
	Appid          string   `xml:"appid,omitempty"`
	Mchid          string   `xml:"mch_id,omitempty"`
	Transaction_id string   `xml:"transaction_id,omitempty"`
	Out_trade_no   string   `xml:"out_trade_no,omitempty"`
	Nonce_str      string   `xml:"nonce_str,omitempty"`
	Sign           string   `xml:"sign,omitempty"`
	Signtype       string   `xml:"sign_type,omitempty"`
}

type OrderqueryRsp struct {
	XMLName      xml.Name `xml:"xml,omitempty"`
	Return_code  string   `xml:"return_code,omitempty"`
	Return_msg   string   `xml:"return_msg,omitempty"`
	Appid        string   `xml:"appid,omitempty"`
	Mch_id       string   `xml:"mch_id,omitempty"`
	Nonce_str    string   `xml:"nonce_str,omitempty"`
	Sign         string   `xml:"sign,omitempty"`
	Result_code  string   `xml:"result_code,omitempty"`
	Err_code     string   `xml:"err_code,omitempty"`
	Err_code_des string   `xml:"err_code_des,omitempty"`
	Device_info  string   `xml:"device_info,omitempty"`
	Openid       string   `xml:"openid,omitempty"`
	Is_subscribe string   `xml:"is_subscribe,omitempty"`
	Trade_type   string   `xml:"trade_type,omitempty"`
	//Trade_state: SUCCESS—支付成功
	//REFUND—转入退款
	//NOTPAY—未支付
	//CLOSED—已关闭
	//REVOKED—已撤销（付款码支付）
	//USERPAYING--用户支付中（付款码支付）
	//PAYERROR--支付失败(其他原因，如银行返回失败)
	Trade_state          string `xml:"trade_state,omitempty"`
	Bank_type            string `xml:"bank_type,omitempty"`
	Fee_type             string `xml:"fee_type,omitempty"`
	Total_fee            int    `xml:"total_fee,omitempty"`
	Settlement_total_fee int    `xml:"settlement_total_fee,omitempty"`
	Coupon_fee           string `xml:"coupon_fee,omitempty"`
	Cash_fee_type        string `xml:"cash_fee_type,omitempty"`
	Cash_fee             string `xml:"cash_fee,omitempty"`
	Transaction_id       string `xml:"transaction_id,omitempty"`
	Out_trade_no         string `xml:"out_trade_no,omitempty"`
	Attach               string `xml:"attach,omitempty"`
	Time_end             string `xml:"time_end,omitempty"`
	Promotion_detail     string `xml:"promotion_detail,omitempty"`
	Trade_state_desc     string `xml:"trade_state_desc,omitempty"` //对当前查询订单状态的描述和下一步操作的指引
}


type Authcodetoopenid struct {
	XMLName        xml.Name `xml:"xml,omitempty"`
	Appid          string   `xml:"appid,omitempty"`
	Mchid          string   `xml:"mch_id,omitempty"`
	Nonce_str      string   `xml:"nonce_str,omitempty"`
	Sign           string   `xml:"sign,omitempty"`
	AuthCode       string   `xml:"auth_code,omitempty"` //付款码
}

type AuthcodetoopenidRsp struct {
	XMLName      xml.Name `xml:"xml,omitempty"`
	Return_code  string   `xml:"return_code,omitempty"`
	Return_msg   string   `xml:"return_msg,omitempty"`
	Appid        string   `xml:"appid,omitempty"`
	Mch_id       string   `xml:"mch_id,omitempty"`
	Nonce_str    string   `xml:"nonce_str,omitempty"`
	Sign         string   `xml:"sign,omitempty"`
	Result_code  string   `xml:"result_code,omitempty"`
	Err_code     string   `xml:"err_code,omitempty"`
	Openid       string   `xml:"openid,omitempty"`
}
