package gongzhonghao

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tang-go/go-dog/lib/encrypt"
	"github.com/tang-go/go-dog/lib/json"
	"github.com/tang-go/go-dog/lib/net"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/weixin/gongzhonghao/handle"
	wxmodel "github.com/tang-go/go-dog/weixin/gongzhonghao/model"
)

type WeixinCfg struct {
	Addr      string
	Appid     string
	Appsecret string
	MchId     string
	APIKey    string
	Token     string
	iplist    []string
	KeyFile   string
	CertFile  string
	WXRoot    string
}

type WeixinCacheInterface interface {
	Get(key string, data interface{}) bool
	Set(key string, data interface{}) bool
	Delete(key string) bool
}

type weixntoken struct {
	Token   string
	TimeOut int64
}

type weixinticket struct {
	Ticket        string
	TicketTimeOut int64
}

//微信主结构
type WeiXinSession struct {
	weixntoken
	weixinticket
	cfg   WeixinCfg
	cache WeixinCacheInterface
}

//在线获取access_token
func (this *WeiXinSession) reloadToken() error {
	//优先缓存获取token
	var token weixntoken
	if this.cache != nil {
		if this.cache.Get(CacheTokenName+this.cfg.Appid, &token) {
			//更新本地内存
			log.Debugln("读取缓存token到本地:", token.Token, " 到期时间:", time.Now().Add(time.Second*time.Duration(token.TimeOut-time.Now().Unix())).String())
			this.weixntoken = token
		}
	}
	log.Debugln("time:", this.TimeOut, "now:", time.Now().Unix(), this.cfg.Appid, this.cfg.Appsecret)
	//还是超时或者没有
	if this.Token == "" || this.TimeOut <= time.Now().Unix() {
		//在线更新token
		var rsp wxmodel.TokenRsp
		r, e := net.HttpsGet(fmt.Sprintf(TokenURL, this.cfg.Appid, this.cfg.Appsecret))
		if e != nil {
			log.Warnln(e)
			return e
		}
		e = json.Unmarshal(r, &rsp)
		if e != nil {
			log.Warnln(e)
			return e
		}
		log.Debugln("在线获取新的token :", rsp.Access_token)
		this.Token = rsp.Access_token
		this.TimeOut = int64(rsp.Expires_in) + time.Now().Unix()
		if rsp.Errcode != 0 {
			return errors.New(rsp.ErrRsp.Errmsg)
		}
		//保存到缓存
		if this.cache != nil {
			this.cache.Set(CacheTokenName+this.cfg.Appid, &this.weixntoken)
		}
	}

	return nil
}

func (this *WeiXinSession) GetAccessToken() (string, error) {
	return this.getToken()
}

//本地获取access_token
func (this *WeiXinSession) getToken() (string, error) {
	//检查token是否已经获取和超时
	if this.Token == "" || this.TimeOut <= time.Now().Unix() {
		if e := this.reloadToken(); e != nil {
			return "", e
		}

	}
	return this.Token, nil
}

//Ticket
func (this *WeiXinSession) reloadTicket() error {
	this.getToken()

	//优先缓存获取Ticket
	var ticket weixinticket
	if this.cache != nil {
		if this.cache.Get(CacheTicketName+this.cfg.Appid, &ticket) {
			//更新本地内存
			log.Debugln("读取缓存ticket到本地:", ticket.Ticket, " 到期时间:", time.Now().Add(time.Second*time.Duration(ticket.TicketTimeOut-time.Now().Unix())).String())
			this.weixinticket = ticket
		}
	}
	log.Debugln("time:", this.TicketTimeOut, "now:", time.Now().Unix())
	//还是超时或者没有
	if this.Ticket == "" || this.TicketTimeOut <= time.Now().Unix() {
		//在线更新token
		var rsp wxmodel.TicketRsp
		r, e := net.HttpsGet(fmt.Sprintf(GetTicketStr, this.Token))
		log.Debugln(fmt.Sprintf(GetTicketStr, this.Token), ": ", string(r))
		if e != nil {
			log.Warnln(e)
			return e
		}
		e = json.Unmarshal(r, &rsp)
		if e != nil {
			log.Warnln(e)
			return e
		}
		log.Debugln("在线获取新的ticket:", rsp.Ticket)
		this.Ticket = rsp.Ticket
		this.TicketTimeOut = int64(rsp.Expires_in) + time.Now().Unix()
		if rsp.Errcode != 0 {
			return errors.New(rsp.ErrRsp.Errmsg)
		}
		//保存到缓存
		if this.cache != nil {
			this.cache.Set(CacheTicketName+this.cfg.Appid, &this.weixinticket)
		}
	}

	return nil
}

//本地获取ticket
func (this *WeiXinSession) getTicket() (string, error) {
	//检查ticket是否已经获取和超时
	if this.Ticket == "" || this.TicketTimeOut <= time.Now().Unix() {
		if e := this.reloadTicket(); e != nil {
			return "", e
		}

	}
	return this.Ticket, nil
}

//获取微信服务器IP地址
func (this *WeiXinSession) getWeixinSvrIp() error {
	//更新token
	var rsp wxmodel.ServerListRsp
	r, e := net.HttpsGet(fmt.Sprintf(GetServerURL, this.Token))
	if e != nil {
		log.Warnln(e)
		return e
	}
	e = json.Unmarshal(r, &rsp)
	if e != nil {
		log.Warnln(e)
		return e
	}
	log.Debugln(rsp)
	if rsp.Errcode != 0 {
		return errors.New(rsp.ErrRsp.Errmsg)
	}
	return nil
}

func (this *WeiXinSession) SetCacheInterface(cache WeixinCacheInterface) {
	this.cache = cache

}

func (this *WeiXinSession) Run(cb handle.WeixinMsgHandleInterface, eventcb handle.WeixinEventHandleInterface) {
	if _, e := this.getToken(); e != nil {
		log.Debugln(e)
	}
	//this.getWeixinSvrIp()
	svrInit(this, cb, eventcb)
}

func (this *WeiXinSession) TokenReset() {
	this.Token = ""
	this.Ticket = ""
	this.cache.Delete(CacheTicketName + this.cfg.Appid)
	this.cache.Delete(CacheTokenName + this.cfg.Appid)
	this.reloadToken()
	this.reloadTicket()
}

func NewSession(cfg WeixinCfg) *WeiXinSession {
	O := new(WeiXinSession)
	O.cfg = cfg
	return O
}

func (this *WeiXinSession) GetUserToken(code string) (error, *wxmodel.UserToken) {
	url := fmt.Sprintf(weixin_user_token, this.cfg.Appid, this.cfg.Appsecret, code)
	body := net.HttpGet(url)
	log.Infoln(string(body))
	usr := new(wxmodel.UserToken)
	if err := json.Unmarshal(body, usr); err != nil {
		return err, nil
	}
	if usr.Errcode != 0 {
		return errors.New(usr.Errmsg), nil
	}
	log.Debugln("GetUserToken %+v", usr)
	return nil, usr
}

func (this *WeiXinSession) GetUInfo(openid, token string) (error, *wxmodel.UInfo) {
	if _, err := this.getToken(); err != nil {
		log.Errorln(err)
		return err, nil
	}

	url := fmt.Sprintf(weixin_user_info, token, openid)
	body := net.HttpGet(url)
	log.Infoln(string(body))
	usr := new(wxmodel.UInfo)
	if err := json.Unmarshal(body, usr); err != nil {
		return err, nil
	}

	log.Debugln("UInfo %+v", usr)
	return nil, usr
}

var TopColor = "#FF0000"

func (this *WeiXinSession) SendTemplateMsg(openid, templateid string, data interface{}) error {
	if _, err := this.getToken(); err != nil {
		log.Errorln(err)
		return err
	}
	wtc := wxmodel.WeixinTemp{}
	wtc.TemplateId = templateid
	wtc.ToUser = openid
	wtc.Topcolor = TopColor
	wtc.Data = data
	buf, _ := json.Marshal(wtc)
	url := weixin_template_url + "?access_token=" + this.Token

	log.Debugln(string(buf))
	body := net.HttpPost(url, []byte(buf))
	wxError := new(wxmodel.ErrRsp)
	if err := json.Unmarshal(body, wxError); err != nil {
		return err
	}

	log.Debugln(wxError)
	if wxError.Errcode != 0 {
		body := net.HttpPost(url, []byte(buf))
		if err := json.Unmarshal(body, wxError); err != nil {
			return err
		}
		if wxError.Errcode != 0 {
			fmt.Println(wxError)
			return errors.New(wxError.Errmsg)
		}
	}
	return nil
}

func (this *WeiXinSession) GetShareSignature(url string) (interface{}, error) {
	count := 0
ag:
	if _, err := this.getTicket(); err != nil {
		this.TokenReset()
		log.Errorln(err)
		//return nil, err
		if count > 10 {
			return nil, err
		}
		count++
		goto ag
	}
	s := new(wxmodel.SignatureResponse)
	s.NonceStr = randStr(16)
	s.Timestamp = time.Now().Unix()
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		this.Ticket, s.NonceStr, s.Timestamp, url)
	log.Debugln("jsapi: ", str)
	h := sha1.New()
	h.Write([]byte(str))
	s.Signature = hex.EncodeToString(h.Sum(nil))
	return s, nil
}

func (this *WeiXinSession) GetUserInfo(openid string) *wxmodel.UserInfoRsp {
	var err error
	if _, err = this.getToken(); err != nil {
		log.Errorln(err)
		return nil
	}
	url := fmt.Sprintf(weixin_userinfo_url, this.Token, openid)

	body := net.HttpGet(url)
	log.Debugln(string(body))
	rsp := new(wxmodel.UserInfoRsp)
	if err := json.Unmarshal(body, rsp); err != nil {
		log.Debugln(err)
		return nil
	}

	return rsp
}

func fillEmpty(s string) string {
	if s == "" {
		return "killEmpty"
	}
	return s
}

const timeFormat = "20060102150405"

//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=4_3
//计算APP签名
//https://pay.weixin.qq.com/wiki/doc/api/app/app.php?chapter=9_1
//body 商品描述
//detail  商品详情
//attach  自定义字段
func (this *WeiXinSession) GetAppUnifiedOrder(openid, orderid, body, detail, notifyurl, attach string, fee int, expire time.Time) (*wxmodel.AppUnifierOrderRsp, error) {
	order := new(wxmodel.AppUnifierOrder)
	order.Appid = this.cfg.Appid
	order.Mchid = this.cfg.MchId
	order.NonceStr = encrypt.GetRandomString(16)
	order.Deviceinfo = "mini"
	order.Body = fillEmpty(body)
	order.Detail = fillEmpty(detail)
	order.Attach = fillEmpty(attach)
	order.Outtradeno = fillEmpty(orderid)
	order.TotalFee = fee
	order.Notifyurl = fillEmpty(notifyurl)
	order.Tradetype = "JSAPI"
	order.Openid = fillEmpty(openid)
	order.Timeexpire = expire.Format(timeFormat)
	stringTemp := fmt.Sprintf("appid=%s&attach=%s&body=%s&detail=%s&device_info=%s&mch_id=%s&nonce_str=%s&notify_url=%s&openid=%s&out_trade_no=%s&time_expire=%s&total_fee=%d&trade_type=%s&key=%s",
		this.cfg.Appid,
		order.Attach,
		order.Body,
		order.Detail,
		order.Deviceinfo,
		this.cfg.MchId,
		order.NonceStr,
		order.Notifyurl,
		order.Openid,
		order.Outtradeno,
		order.Timeexpire,
		order.TotalFee,
		order.Tradetype,
		this.cfg.APIKey,
	)
	log.Infoln(stringTemp)
	order.Sign = strings.ToUpper(encrypt.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Debugln("--GetAppUnifiedOrder--", string(buf))
	r := net.HttpPost(app_pay_unifiedorder, []byte(buf))
	log.Debugln("--Rsp--", string(r))
	rsp := new(wxmodel.AppUnifierOrderRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Debugln("--AppUnifierOrderRsp--", rsp)
	return rsp, nil
}

//扫码支付
func (this *WeiXinSession) Qrpay(authcode, orderid, body, detail, attach string, fee int) (*wxmodel.QrPayRsp, error) {
	order := new(wxmodel.QrPay)
	order.Appid = this.cfg.Appid
	order.Mchid = this.cfg.MchId
	order.NonceStr = encrypt.GetRandomString(16)
	order.Deviceinfo = "qrcode"
	order.Body = fillEmpty(body)
	order.Detail = fillEmpty(detail)
	order.Attach = fillEmpty(attach)
	order.Outtradeno = fillEmpty(orderid)
	order.TotalFee = fee
	order.AuthCode = authcode
	//order.TimeStart = time.Now().Format(timeFormat)
	//order.Timeexpire = time.Now().Add(time.Minute).Format(timeFormat)
	//stringTemp := fmt.Sprintf("appid=%s&attach=%s&body=%s&detail=%s&device_info=%s&mch_id=%s&nonce_str=%s&out_trade_no=%s&time_expire=%s&time_start=%s&total_fee=%d&key=%s",
	//	this.cfg.Appid,
	//	order.Attach,
	//	order.Body,
	//	order.Detail,
	//	order.Deviceinfo,
	//	this.cfg.MchId,
	//	order.NonceStr,
	//	order.Outtradeno,
	//	order.Timeexpire,
	//	order.TimeStart,
	//	order.TotalFee,
	//	this.cfg.APIKey,
	//)
	stringTemp := fmt.Sprintf("appid=%s&attach=%s&auth_code=%s&body=%s&detail=%s&device_info=%s&mch_id=%s&nonce_str=%s&out_trade_no=%s&total_fee=%d&key=%s",
		this.cfg.Appid,
		order.Attach,
		order.AuthCode,
		order.Body,
		order.Detail,
		order.Deviceinfo,
		order.Mchid,
		order.NonceStr,
		order.Outtradeno,
		order.TotalFee,
		this.cfg.APIKey,
	)
	log.Infoln(stringTemp)
	order.Sign = strings.ToUpper(encrypt.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Debugln("--app_pay_qrcode--", string(buf))
	r := net.HttpPost(app_pay_qrcode, []byte(buf))
	log.Debugln("--Rsp--", string(r))
	rsp := new(wxmodel.QrPayRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Debugln("--AppUnifierOrderRsp--", rsp)
	return rsp, nil
}

//查询订单状态
func (this *WeiXinSession) Orderquery(transaction_id, out_trade_no string) (*wxmodel.OrderqueryRsp, error) {
	var stringTemp string
	order := new(wxmodel.Orderquery)
	order.Appid = this.cfg.Appid
	order.Mchid = this.cfg.MchId
	order.Transaction_id = transaction_id
	order.Out_trade_no = out_trade_no
	order.Nonce_str = encrypt.GetRandomString(16)
	if order.Out_trade_no != "" {
		stringTemp = fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&out_trade_no=%s&key=%s",
			this.cfg.Appid,
			order.Mchid,
			order.Nonce_str,
			order.Out_trade_no,
			this.cfg.APIKey,
		)
	} else {
		stringTemp = fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&transaction_id=%s&key=%s",
			this.cfg.Appid,
			order.Mchid,
			order.Nonce_str,
			order.Transaction_id,
			this.cfg.APIKey,
		)
	}
	log.Debugln(stringTemp)
	order.Sign = strings.ToUpper(encrypt.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	r := net.HttpPost(app_pay_orderquery, []byte(buf))
	log.Debugln("--Rsp--", string(r))
	rsp := new(wxmodel.OrderqueryRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Debugln("--OrderqueryRsp--", rsp)
	return rsp, nil
}

// 付款码查询openid
func (this *WeiXinSession) Authcodetoopenid(AuthCode string) (*wxmodel.AuthcodetoopenidRsp, error) {
	var stringTemp string
	order := new(wxmodel.Authcodetoopenid)
	order.Appid = this.cfg.Appid
	order.Mchid = this.cfg.MchId
	order.Nonce_str = encrypt.GetRandomString(16)
	order.AuthCode = AuthCode
	stringTemp = fmt.Sprintf("appid=%s&auth_code=%s&mch_id=%s&nonce_str=%s&key=%s",
		this.cfg.Appid,
		order.AuthCode,
		order.Mchid,
		order.Nonce_str,
		this.cfg.APIKey,
	)
	log.Debugln(stringTemp)
	order.Sign = strings.ToUpper(encrypt.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	r := net.HttpPost(app_pay_authcodetoopenid, []byte(buf))
	log.Debugln("--Rsp--", string(r))
	rsp := new(wxmodel.AuthcodetoopenidRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Debugln("--AuthcodetoopenidRsp--", rsp)
	return rsp, nil
}

func (this *WeiXinSession) Refund(transaction_id, notifyurl, out_refund_order, desc string, fee, total int) (*wxmodel.RefundRsp, error) {
	order := new(wxmodel.RefundReq)
	order.Appid = this.cfg.Appid
	order.Mchid = this.cfg.MchId
	order.NonceStr = encrypt.GetRandomString(16)
	order.TotalFee = total
	order.RefundFee = fee
	order.Transaction_id = transaction_id
	order.Notifyurl = fillEmpty(notifyurl)
	order.Refund_desc = desc
	order.Out_refund_no = out_refund_order
	stringTemp := fmt.Sprintf("appid=%s&mch_id=%s&nonce_str=%s&notify_url=%s&out_refund_no=%s&refund_desc=%s&refund_fee=%d&total_fee=%d&transaction_id=%s&key=%s",
		this.cfg.Appid,
		this.cfg.MchId,
		order.NonceStr,
		order.Notifyurl,
		order.Out_refund_no,
		order.Refund_desc,
		order.RefundFee,
		order.TotalFee,
		order.Transaction_id,
		this.cfg.APIKey,
	)
	log.Infoln(stringTemp)
	order.Sign = strings.ToUpper(encrypt.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Debugln("--Refund--", string(buf), "--", this.cfg.KeyFile, "--", this.cfg.CertFile)
	r := net.HttpsPost(refund_url, this.cfg.KeyFile, this.cfg.CertFile, this.cfg.WXRoot, []byte(buf))
	log.Debugln("--Rsp--", string(r))
	rsp := new(wxmodel.RefundRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Debugln("--Refund--", rsp)
	return rsp, nil
}

func (gzh *WeiXinSession) WxConfig(url string) (conf wxmodel.WxConfig) {
	count := 0
ag:
	if _, err := gzh.getTicket(); err != nil {
		gzh.TokenReset()
		log.Errorln(err)
		//return nil, err
		if count > 10 {
			return conf
		}
		count++
		goto ag
	}

	conf.AppId = gzh.cfg.Appid
	conf.Timestamp = time.Now().Unix()
	conf.NonceStr = randStr(16)
	sh := sha1.New()
	s := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		gzh.Ticket,
		conf.NonceStr,
		conf.Timestamp,
		url,
	)
	log.Infoln(s)
	sh.Write([]byte(s))
	conf.Signature = fmt.Sprintf("%x", sh.Sum(nil))
	return conf
}

func (gzh *WeiXinSession) Draw(orderno, openid, desc, ip string, amount int) (*wxmodel.UserDrawRsp, error) {
	order := new(wxmodel.UserDraw)
	order.Mchid = gzh.cfg.MchId
	order.Nonce_str = encrypt.GetRandomString(16)
	order.Partner_trade_no = orderno
	order.Mch_appid = gzh.cfg.Appid
	order.Openid = openid
	order.Check_name = "NO_CHECK"
	order.Amount = amount
	order.Desc = desc
	order.Spbill_create_ip = ip
	stringTemp := fmt.Sprintf("amount=%d&check_name=%s&desc=%s&mch_appid=%s&mchid=%s&nonce_str=%s&openid=%s&partner_trade_no=%s&spbill_create_ip=%s&key=%s",
		order.Amount,
		order.Check_name,
		order.Desc,
		order.Mch_appid,
		order.Mchid,
		order.Nonce_str,
		order.Openid,
		order.Partner_trade_no,
		order.Spbill_create_ip,
		gzh.cfg.APIKey,
	)
	log.Infoln(stringTemp)
	order.Sign = strings.ToUpper(encrypt.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Debugln("--draw--", string(buf))
	r := net.HttpsPost(draw_ulr, gzh.cfg.KeyFile, gzh.cfg.CertFile, gzh.cfg.WXRoot, []byte(buf))
	log.Debugln("--draw--", string(r))
	rsp := new(wxmodel.UserDrawRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Debugln("--draw--", rsp)
	return rsp, nil
}
