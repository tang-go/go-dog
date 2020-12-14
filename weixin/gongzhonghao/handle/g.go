package handle

import (
	"encoding/xml"

	"github.com/tang-go/go-dog/log"
	wxmodel "github.com/tang-go/go-dog/weixin/gongzhonghao/model"
)

type WeixinMsgHandleInterface interface {
	OnText(msg *wxmodel.WexinTextMessage) interface{}
	OnImage(msg *wxmodel.WexinImageMessage) interface{}
	OnVoice(msg *wxmodel.WexinVoiceMessage) interface{}
	OnVideo(msg *wxmodel.WexinVideoMessage) interface{}
	OnLink(msg *wxmodel.WexinLinkMessage) interface{}
	OnLocation(msg *wxmodel.WexinLocationMessage) interface{}
}

var weixinMsgHandleImpl WeixinMsgHandleInterface
var weixinEventHanldeImpl WeixinEventHandleInterface

type WeixinEventHandleInterface interface {
	OnSubscribe(event *wxmodel.WeixinSubscribeEvent) interface{}
	OnScanner(event *wxmodel.WeixinScannerEvent) interface{}
	OnLocationEvent(event *wxmodel.WeixinLocationEvent) interface{}
	OnMenu(event *wxmodel.WeixinMenuEvent) interface{}
}

type handle interface {
	GetKey() string
	Handle(data []byte) interface{}
}

var handles map[string]handle
var eventhandles map[string]handle

//处理主入口
func Entrance(head wxmodel.WeixinBase, data []byte) interface{} {
	if head.MsgType.Value == "event" {
		var msg wxmodel.WeixinEventBase
		if e := xml.Unmarshal(data, &msg); e != nil {
			log.Warnln(e)
		}
		if h, ok := eventhandles[msg.Event.Value]; ok {
			return h.Handle(data)
		} else {
			log.Warnln(msg.Event.Value, " event not found  ")
		}
	} else {
		if h, ok := handles[head.MsgType.Value]; ok {
			return h.Handle(data)
		} else {
			log.Warnln(head.MsgType.Value, " not found  ")
		}
	}
	return nil
}

func register(h handle) {
	handles[h.GetKey()] = h
}
func registerevent(h handle) {
	eventhandles[h.GetKey()] = h
}
func Init(cb WeixinMsgHandleInterface, event WeixinEventHandleInterface) {
	weixinMsgHandleImpl = cb
	weixinEventHanldeImpl = event
	handles = make(map[string]handle, 0)
	eventhandles = make(map[string]handle, 0)

	register(new(TextHandle))
	register(new(ImageHandle))
	register(new(VoiceHandle))
	register(new(LocationHandle))
	register(new(VideoHandle))
	register(new(ShortVideoHandle))
	register(new(LinkHandle))
	registerevent(new(SubscribeEventHandle))
	registerevent(new(ScannerEventHandle))
	registerevent(new(LocationEventHandle))
	registerevent(new(MenuEventHandle))
}
