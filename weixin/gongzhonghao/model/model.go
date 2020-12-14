package model

import "encoding/xml"

type ErrRsp struct {
	Errcode int
	Errmsg  string
}

type TokenRsp struct {
	ErrRsp
	Access_token string
	Expires_in   int
}

type TicketRsp struct {
	ErrRsp
	Ticket     string
	Expires_in int
}

type SignatureResponse struct {
	NonceStr  string `json:"nonceStr"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

type ServerListRsp struct {
	ErrRsp
	Ip_list []string
}

type CdataString struct {
	Value string `xml:",cdata"`
}

//微信接受消息的基础参数（每一种类型的消息都需要这些参数）
type WeixinBase struct {
	ToUserName   CdataString `xml:"ToUserName"`
	FromUserName CdataString `xml:"FromUserName"`
	CreateTime   int64       `xml:"CreateTime"`
	MsgType      CdataString `xml:"MsgType"`
}

type WeixinTextBase struct {
	WeixinBase
	MsgId int64 `xml:"MsgId"`
}

//接收微信文本消息
type WexinTextMessage struct {
	XMLName xml.Name `xml:"xml"`
	WeixinTextBase
	Content CdataString `xml:"Content"`
}

func (this *WexinTextMessageRsp) PrepareWeixinTextMessage(base WeixinBase, text CdataString) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.Content = text
}

type WexinTextMessageRsp WexinTextMessage

//接收微信图片消息
type WexinImageMessage struct {
	XMLName xml.Name `xml:"xml"`
	WeixinBase
	PicUrl  string      `xml:"PicUrl"`
	MediaId CdataString `xml:"MediaId"`
}

func (this *WexinImageMessageRsp) PrepareWexinImageMessage(base WeixinBase, url string, mediaId CdataString) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.PicUrl = url
	this.MediaId = mediaId
}

type WexinImageMessageRsp WexinImageMessage

//接收微信语音消息
type WexinVoiceMessage struct {
	XMLName xml.Name `xml:"xml"`
	WeixinBase
	Format      string      `xml:"Format"`
	MediaId     CdataString `xml:"MediaId"`
	Recognition string      `xml:"Recognition"`
}

func (this *WexinVoiceMessageRsp) PrepareWexinImageMessage(base WeixinBase, format string, mediaId CdataString, recognition string) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.Format = format
	this.MediaId = mediaId
	this.Recognition = recognition
}

type WexinVoiceMessageRsp WexinVoiceMessage

//接收微信视频消息
type WexinVideoMessage struct {
	XMLName xml.Name `xml:"xml"`
	WeixinBase
	ThumbMediaId string `xml:"ThumbMediaId"`
	MediaId      string `xml:"MediaId"`
}

func (this *WexinVideoMessageRsp) PrepareWexinVideoMessage(base WeixinBase, thumbMediaId string, mediaId string) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.ThumbMediaId = thumbMediaId
	this.MediaId = mediaId
}

type WexinVideoMessageRsp WexinVideoMessage

type WexinShortVideoMessage WexinVideoMessage

//地理位置消息
type WexinLocationMessage struct {
	XMLName xml.Name `xml:"xml"`
	WeixinBase
	Location_X float32 `xml:"Location_X"`
	Location_Y float32 `xml:"Location_Y"`
	Scale      int     `xml:"Scale"`
	Label      string  `xml:"Label"`
}

func (this *WexinLocationMessageRsp) PrepareWexinLocationMessage(base WeixinBase, location_X float32, location_Y float32, scale int, label string) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.Location_X = location_X
	this.Location_Y = location_Y
	this.Scale = scale
	this.Label = label
}

type WexinLocationMessageRsp WexinLocationMessage

//链接消息
type WexinLinkMessage struct {
	XMLName xml.Name `xml:"xml"`
	WeixinBase
	Title       string `xml:"Title"`
	Description string `xml:"Description"`
	Url         string `xml:"Url"`
}

func (this *WexinLinkMessageRsp) PrepareWexinLinkMessage(base WeixinBase, title string, description string, url string) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.Title = title
	this.Description = description
	this.Url = url
}

type WexinLinkMessageRsp WexinLinkMessage

//接收事件推送的基础结构
type WeixinEventBase struct {
	WeixinBase
	Event CdataString `xml:"Event"` //事件类型
}

//关注/取消关注事件
type WeixinSubscribeEvent struct {
	XMLName xml.Name `xml:"xml"`
	WeixinEventBase
}

func (this *WeixinSubscribeEventRsp) PrepareWeixinSubscribeEvent(base WeixinEventBase) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.Event = base.Event
}

type WeixinSubscribeEventRsp WeixinSubscribeEvent

//扫描带参数二维码事件
type WeixinScannerEvent struct {
	XMLName xml.Name `xml:"xml"`
	WeixinEventBase
	EventKey string `xml:"EventKey"` //事件KEY值，qrscene_为前缀，后面为二维码的参数值
	Ticket   string `xml:"Ticket"`   //二维码的ticket，可用来换取二维码图片
}

func (this *WeixinScannerEventRsp) PrepareWeixinScannerEvent(base WeixinEventBase, eventKey string, ticket string) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.Event = base.Event
	this.EventKey = eventKey
	this.Ticket = ticket
}

type WeixinScannerEventRsp WeixinScannerEvent

//上报地理位置事件
type WeixinLocationEvent struct {
	XMLName xml.Name `xml:"xml"`
	WeixinEventBase
	Latitude  float32 `xml:"Latitude"`  //地理位置纬度
	Longitude float32 `xml:"Longitude"` //地理位置经度
	Precision float32 `xml:"Precision"` //地理位置精度
}

func (this *WeixinLocationEventRsp) PrepareWeixinLocationEvent(base WeixinEventBase, latitude float32, longitude float32, precision float32) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.Event = base.Event
	this.Latitude = latitude
	this.Longitude = longitude
	this.Precision = precision
}

type WeixinLocationEventRsp WeixinLocationEvent

//自定义菜单事件
type WeixinMenuEvent struct {
	XMLName xml.Name `xml:"xml"`
	WeixinEventBase
	EventKey string `xml:"EventKey"` //事件KEY值，与自定义菜单接口中KEY值对应
}

func (this *WeixinMenuEventRsp) PrepareWeixinMenuEvent(base WeixinEventBase, eventKey string) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.Event = base.Event
	this.EventKey = eventKey
}

type WeixinMenuEventRsp WeixinMenuEvent

type UserInfoRsp struct {
	ErrRsp
	Subscribe       int    `json:"subscribe"`
	Openid          string `json:"openid"`
	Nickname        string `json:"nickname"`
	Sex             int    `json:"sex"`
	Language        string `json:"language"`
	City            string `json:"city"`
	Province        string `json:"province"`
	Country         string `json:"country"`
	Headimgurl      string `json:"headimgurl"`
	Subscribe_time  int64  `json:"subscribe_time"`
	Unionid         string `json:"unionid"`
	Remark          string `json:"remark"`
	Groupid         int    `json:"groupid"`
	Tagid_list      []int  `json:"tagid_list"`
	Subscribe_scene string `json:"subscribe_scene"`
	Qr_scene        int    `json:"qr_scene"`
	Qr_scene_str    string `json:"qr_scene_str"`
}

type UserToken struct {
	Access_token string `json:"access_token"`
	Expires_in   int    `json:"expires_in"`
	Openid       string `json:"openid"`
	Scope        string `json:"scope"`
	Unionid      string `json:"unionid"`
	ErrRsp
}

type UInfo struct {
	Openid     string   `json:"openid"`
	Nickname   string   `json:"nickname"`
	Sex        int      `json:"sex"`
	Province   string   `json:"province"`
	City       string   `json:"city"`
	Country    string   `json:"country"`
	Headimgurl string   `json:"headimgurl"`
	Privilege  []string `json:"privilege"`
	Unionid    string   `json:"unionid"`
}
