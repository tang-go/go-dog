package xiaochengxu

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tang-go/go-dog/lib/encrypt"
	"github.com/tang-go/go-dog/lib/json"
	"github.com/tang-go/go-dog/lib/net"
	"github.com/tang-go/go-dog/lib/utils"
	wxmodel "github.com/tang-go/go-dog/weixin/gongzhonghao/model"
	"github.com/tang-go/go-dog/weixin/xiaochengxu/model"
	minimodel "github.com/tang-go/go-dog/weixin/xiaochengxu/model"

	"github.com/tang-go/go-dog/log"
)

type WeixinCfg struct {
	Appid     string
	Appsecret string
	MchId     string
	APIKey    string
	KeyFile   string
	CertFile  string
	WXRoot    string
}

type MiniTpl struct {
	Tplid   string
	Context string
}

type weixntoken struct {
	Token   string
	TimeOut int64
}

type WeixinCacheInterface interface {
	Get(key string, data interface{}) bool
	Set(key string, data interface{}) bool
	Delete(key string) bool
}

//微信小程序主结构
type WeiXinMiniSession struct {
	weixntoken
	cfg   WeixinCfg
	cache WeixinCacheInterface
}

//在线获取access_token
func (this *WeiXinMiniSession) reloadToken() error {
	//优先缓存获取token
	var token weixntoken
	if this.cache != nil {
		if this.cache.Get(CacheTokenName+this.cfg.Appid, &token) {
			//更新本地内存
			log.Debugln("读取缓存token到本地:", token.Token, " 到期时间:", time.Now().Add(time.Second*time.Duration(token.TimeOut-time.Now().Unix())).String())
			this.weixntoken = token
		}
	}
	log.Debugln("time:", this.TimeOut, "now:", time.Now().Unix())
	//还是超时或者没有
	if this.Token == "" || this.TimeOut <= time.Now().Unix() {
		//在线更新token
		var rsp model.AccessTokenRsp
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
		this.Token = rsp.Access_token
		this.TimeOut = int64(rsp.Expires_in) + time.Now().Unix()
		log.Debugln("在线获取新的token :", rsp.Access_token, " 超时", time.Unix(this.TimeOut, 0))

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

func (this *WeiXinMiniSession) GetAccessToken() (string, error) {
	return this.getToken()
}

//本地获取access_token
func (this *WeiXinMiniSession) getToken() (string, error) {
	//检查token是否已经获取和超时
	if this.Token == "" || this.TimeOut <= time.Now().Unix() {
		if e := this.reloadToken(); e != nil {
			return "", e
		}

	}
	return this.Token, nil
}

func (this *WeiXinMiniSession) SetCacheInterface(cache WeixinCacheInterface) {
	this.cache = cache
}

func (this *WeiXinMiniSession) Run() {
	if _, e := this.getToken(); e != nil {
		log.Debugln(e)
	}
}

func NewSession(cfg WeixinCfg) *WeiXinMiniSession {
	O := new(WeiXinMiniSession)
	O.cfg = cfg
	return O
}

type WxMiniBody struct {
	OpenId      string `json:"openId"`
	Session_key string `json:"session_key"`
	UnionId     string `json:"unionid"`
}

func (sess *WeiXinMiniSession) weixinMiniProgramLogin(code string) string {
	body, err := net.HttpsGet("https://api.weixin.qq.com/sns/jscode2session?appid=" + sess.cfg.Appid + "&secret=" + sess.cfg.Appsecret + "&js_code=" + code + "&grant_type=authorization_code")
	if err != nil {
		log.Debugln(err, "--------->获取openid和sessionKey失败")
		return ""
	}
	pre := string(body)
	if strings.Contains(pre, "errcode") {
		log.Warnln(string(body))
		return ""
	}
	log.Debugln("wx resp", string(body))
	o := new(WxMiniBody)
	json.Unmarshal(body, o)
	return o.Session_key
}

func (sess *WeiXinMiniSession) WxLogin(encryptedData, iv, code string) (*minimodel.UserWX, error) {
	sessionKey := sess.weixinMiniProgramLogin(code)
	log.Debugln("iv:", iv, " code:", code, " session_key:", sessionKey, "encryptedData:", encryptedData)
	pc := utils.WxBizDataCrypt{AppID: sess.cfg.Appid, SessionKey: sessionKey}
	str, err := pc.Decrypt(encryptedData, iv, true) //第三个参数解释： 需要返回 JSON 数据类型时 使用 true, 需要返回 map 数据类型时 使用 false
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	user := new(minimodel.UserWX)
	json.Unmarshal([]byte(str.(string)), user)
	log.Debugln(str.(string))
	log.Debugln(user)
	return user, nil
}

func (sess *WeiXinMiniSession) WxPhone(encryptedData, iv, code string) (*minimodel.UserPhone, error) {
	sessionKey := sess.weixinMiniProgramLogin(code)
	log.Debugln("iv:", iv, " code:", code, " session_key:", sessionKey)
	pc := utils.WxBizDataCrypt{AppID: sess.cfg.Appid, SessionKey: sessionKey}
	str, err := pc.Decrypt(encryptedData, iv, true) //第三个参数解释： 需要返回 JSON 数据类型时 使用 true, 需要返回 map 数据类型时 使用 false
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	user := new(minimodel.UserPhone)
	json.Unmarshal([]byte(str.(string)), user)
	log.Debugln(str.(string))
	log.Debugln(user)
	return user, nil
}

func (this *WeiXinMiniSession) SubscribeMessage(openid, templateid, page string, data interface{}) error {
	var err error
	if _, err = this.getToken(); err != nil {
		log.Errorln(err)
		return err
	}
	wtc := new(model.SubscribeMessage)
	wtc.Template_id = templateid
	wtc.Touser = openid
	wtc.Data = data
	wtc.Miniprogram_state = "developer"
	wtc.Page = page
	buf, _ := json.Marshal(wtc)
	url := mini_subscribe_url + "?access_token=" + this.Token

	log.Debugln(string(buf))
	body := net.HttpPost(url, []byte(buf))
	wxError := new(model.ErrRsp)
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

func computeHmacSha256(message string, secret string) string {
	key := []byte(secret)
	h := hmac.New(sha256.New, key)
	h.Write([]byte(message))
	//	fmt.Println(h.Sum(nil))
	sha := hex.EncodeToString(h.Sum(nil))
	//	fmt.Println(sha)

	//	hex.EncodeToString(h.Sum(nil))
	return base64.StdEncoding.EncodeToString([]byte(sha))
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
func (this *WeiXinMiniSession) GetAppUnifiedOrder(openid, orderid, body, detail, notifyurl, attach string, fee int, expire time.Time) (*wxmodel.AppUnifierOrderRsp, error) {
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
func (this *WeiXinMiniSession) Refund(transaction_id, notifyurl, out_refund_order, desc string, fee, total int) (*wxmodel.RefundRsp, error) {
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
func (this *WeiXinMiniSession) Draw(orderno, openid, desc, ip string, amount int) (*wxmodel.UserDrawRsp, error) {
	order := new(wxmodel.UserDraw)
	order.Mchid = this.cfg.MchId
	order.Nonce_str = encrypt.GetRandomString(16)
	order.Partner_trade_no = orderno
	order.Mch_appid = this.cfg.Appid
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
		this.cfg.APIKey,
	)
	log.Infoln(stringTemp)
	order.Sign = strings.ToUpper(encrypt.MD5(stringTemp))
	buf, err := xml.Marshal(order)
	if err != nil {
		log.Errorln(err)
		return nil, err
	}
	log.Debugln("--draw--", string(buf))
	r := net.HttpsPost(draw_ulr, this.cfg.KeyFile, this.cfg.CertFile, this.cfg.WXRoot, []byte(buf))
	log.Debugln("--draw--", string(r))
	rsp := new(wxmodel.UserDrawRsp)
	if err := xml.Unmarshal(r, rsp); err != nil {
		return nil, err
	}
	log.Debugln("--draw--", rsp)
	return rsp, nil
}
