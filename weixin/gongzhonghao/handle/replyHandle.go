package handle

import (
	"encoding/xml"

	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/weixin/gongzhonghao/model"
)

//回复文本消息
type ReplyTextHandle struct {
}

func (this *ReplyTextHandle) GetKey() string {
	return "回复文本消息"
}

func (this *ReplyTextHandle) ReplyHandle(data []byte) interface{} {
	msg := new(model.ReplyTextMessage)
	e := xml.Unmarshal(data, msg)
	if e != nil {
		log.Warnln(e)
		return nil
	}
	if WeixinReplyHandleImpl != nil {
		return WeixinReplyHandleImpl.OnTextMsg(msg)
	}
	return nil
}

//回复图片消息
type ReplyImageHandle struct {
}

func (this *ReplyImageHandle) GetKey() string {
	return "回复图片消息"
}
func (this *ReplyImageHandle) ReplyHandle(data []byte) interface{} {
	msg := new(model.ReplyImageMessage)
	e := xml.Unmarshal(data, msg)
	if e != nil {
		log.Warnln(e)
		return nil
	}
	if WeixinReplyHandleImpl != nil {
		return WeixinReplyHandleImpl.OnImageMsg(msg)
	}
	return nil
}

//回复语言消息
type ReplyVoicehandle struct {
}

func (this *ReplyVoicehandle) GetKey() string {
	return "回复语音消息"
}
func (this *ReplyVoicehandle) ReplyHandle(data []byte) interface{} {
	msg := new(model.ReplyVoiceMessage)
	e := xml.Unmarshal(data, msg)
	if e != nil {
		log.Warnln(e)
		return nil
	}
	if WeixinReplyHandleImpl != nil {
		return WeixinReplyHandleImpl.OnVoiceMsg(msg)
	}
	return nil
}

//回复视频消息
type ReplyVideoHandle struct {
}

func (this *ReplyVideoHandle) GetKey() string {
	return "回复视频消息"
}
func (this *ReplyVideoHandle) ReplyHandle(data []byte) interface{} {
	msg := new(model.ReplyVideoMessage)
	e := xml.Unmarshal(data, msg)
	if e != nil {
		log.Warnln(e)
		return nil
	}
	if WeixinReplyHandleImpl != nil {
		return WeixinReplyHandleImpl.OnVideoMsg(msg)
	}
	return nil
}

//回复音乐消息
type ReplyMusicHandle struct {
}

func (this *ReplyMusicHandle) GetKey() string {
	return "回复音乐消息"
}
func (this *ReplyMusicHandle) ReplyHandle(data []byte) interface{} {
	msg := new(model.ReplyMusicMessage)
	e := xml.Unmarshal(data, msg)
	if e != nil {
		log.Warnln(e)
		return nil
	}
	if WeixinReplyHandleImpl != nil {
		return WeixinReplyHandleImpl.OnMusicMsg(msg)
	}
	return nil
}

//回复图文消息
type ReplyImageTextHandle struct {
}

func (this *ReplyImageTextHandle) GetKey() string {
	return "回复图文消息"
}
func (this *ReplyImageTextHandle) ReplyHandle(data []byte) interface{} {
	msg := new(model.ReplyImageTextMessage)
	e := xml.Unmarshal(data, msg)
	if e != nil {
		log.Warnln(e)
		return nil
	}
	if WeixinReplyHandleImpl != nil {
		return WeixinReplyHandleImpl.OnImageTextMsg(msg)
	}
	return nil
}
