package model

import (
	"encoding/xml"
	"time"
)

type ReplyMessageBase struct {
	ToUserName   CdataString `xml:"ToUserName"`   //接收方帐号（收到的OpenID）
	FromUserName CdataString `xml:"FromUserName"` //开发者微信号
	CreateTime   int64       `xml:"CreateTime"`   //消息创建时间
	MsgType      CdataString `xml:"MsgType"`      //消息类型
}

//回复文本信息
type ReplyTextMessage struct {
	XMLName xml.Name `xml:"xml"`
	ReplyMessageBase
	Content CdataString `xml:"Content"` //回复的消息内容（换行：在content中能够换行，微信客户端就支持换行显示）
}

func CreateTextMsg(base WeixinBase, content CdataString) *ReplyTextMessage {
	msg := new(ReplyTextMessage)
	msg.ToUserName = base.FromUserName
	msg.FromUserName = base.ToUserName
	msg.CreateTime = time.Now().Unix()
	msg.MsgType = CdataString{Value: "text"}
	msg.Content = content
	return msg
}

//回复图片信息
type ReplyImageMessage struct {
	XMLName xml.Name `xml:"xml"`
	ReplyMessageBase
	MediaId CdataString `xml:"MediaId"` //通过素材管理中的接口上传多媒体文件，得到的id。
}

func CreateImageMsg(base WeixinBase, mediaId CdataString) *ReplyImageMessage {
	msg := new(ReplyImageMessage)
	msg.ToUserName = base.FromUserName
	msg.FromUserName = base.ToUserName
	msg.CreateTime = time.Now().Unix()
	msg.MsgType = CdataString{Value: "image"}
	msg.MediaId = mediaId
	return msg
}

//回复语音消息
type ReplyVoiceMessage struct {
	XMLName xml.Name `xml:"xml"`
	ReplyMessageBase
	MediaId CdataString `xml:"MediaId"` //通过素材管理中的接口上传多媒体文件，得到的id
}

func (this *ReplyVoiceMessageRsp) SendVoiceMsg(base ReplyMessageBase, mediaId CdataString) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.MediaId = mediaId
}

type ReplyVoiceMessageRsp ReplyVoiceMessage

//回复视频消息
type ReplyVideoMessage struct {
	XMLName xml.Name `xml:"xml"`
	ReplyMessageBase
	MediaId     CdataString `xml:"MediaId"`     //通过素材管理中的接口上传多媒体文件，得到的id
	Title       CdataString `xml:"Title"`       //视频消息的标题 非必填
	Description CdataString `xml:"Description"` //视频消息的描述 非必填
}

func (this *ReplyVideoMessageRsp) SendVideoMsg(base ReplyMessageBase, mediaId CdataString, title CdataString, description CdataString) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.MediaId = mediaId
	this.Title = title
	this.Description = description
}

type ReplyVideoMessageRsp ReplyVideoMessage

//回复音乐消息
type ReplyMusicMessage struct {
	XMLName xml.Name `xml:"xml"`
	ReplyMessageBase
	ThumbMediaId CdataString `xml:"ThumbMediaId"` //缩略图的媒体id，通过素材管理中的接口上传多媒体文件，得到的id
	Title        CdataString `xml:"Title"`        //音乐标题 非必填
	Description  CdataString `xml:"Description"`  //音乐描述 非必填
	MusicURL     CdataString `xml:"MusicURL"`     //音乐链接 非必填
	HQMusicUrl   CdataString `xml:"HQMusicUrl"`   //高质量音乐链接，WIFI环境优先使用该链接播放音乐 非必填
}

func (this *ReplyMusicMessageRsp) SendMusicMsg(base ReplyMessageBase, thumbMediaId CdataString, title CdataString, description CdataString, MusicURL CdataString, HQMusicUrl CdataString) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.ThumbMediaId = thumbMediaId
	this.Title = title
	this.Description = description
	this.MusicURL = MusicURL
	this.HQMusicUrl = HQMusicUrl
}

type ReplyMusicMessageRsp ReplyMusicMessage

type ArticleItem struct {
	Title       CdataString `xml:"Title"`       //图文消息标题
	Description CdataString `xml:"Description"` //图文消息描述
	PicUrl      CdataString `xml:"PicUrl"`      //图片链接，支持JPG、PNG格式，较好的效果为大图360*200，小图200*200
	Url         CdataString `xml:"Url"`         //点击图文消息跳转链接
}
type Article struct {
	Items []ArticleItem `xml:"item"` //多条图文消息信息，默认第一个item为大图,注意，如果图文数超过8，则将会无响应
}

//回复图文消息
type ReplyImageTextMessage struct {
	XMLName xml.Name `xml:"xml"`
	ReplyMessageBase
	ArticleCount int64   `xml:"ArticleCount"` //图文消息个数，限制为8条以内
	Articles     Article `xml:"Articles"`     //多条图文消息信息，默认第一个item为大图,注意，如果图文数超过8，则将会无响应

}

func (this *ReplyImageTextMessageRsp) SendImageTextMsg(base ReplyMessageBase, articleCount int64, title CdataString, description CdataString, picUrl CdataString, url CdataString) {
	this.ToUserName = base.ToUserName
	this.FromUserName = base.FromUserName
	this.CreateTime = base.CreateTime
	this.MsgType = base.MsgType
	this.ArticleCount = articleCount

}

type ReplyImageTextMessageRsp ReplyImageTextMessage

func CreateImageTextMessage(base WeixinBase, articleCount int64, title CdataString, description CdataString, picUrl CdataString, url CdataString) *ReplyImageTextMessage {
	msg := new(ReplyImageTextMessage)
	msg.ToUserName = base.FromUserName
	msg.FromUserName = base.ToUserName
	msg.CreateTime = time.Now().Unix()
	msg.MsgType = CdataString{Value: "news"}
	msg.ArticleCount = 1
	msg.Articles.Items = append(msg.Articles.Items, ArticleItem{Title: title, Url: url, PicUrl: picUrl, Description: description})
	return msg
}
