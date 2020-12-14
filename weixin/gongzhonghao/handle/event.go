package handle

import (
	"encoding/xml"

	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/weixin/gongzhonghao/model"
)

//处理关注和取消关注事件
type SubscribeEventHandle struct {
}

func (this *SubscribeEventHandle) GetKey() string {
	return "subscribe"
}

func (this *SubscribeEventHandle) Handle(data []byte) interface{} {
	event := new(model.WeixinSubscribeEvent)

	e := xml.Unmarshal(data, event)

	if e != nil {
		log.Warnln(e)
		return nil
	}

	if weixinEventHanldeImpl != nil {
		return weixinEventHanldeImpl.OnSubscribe(event)
	}

	return nil
}

//处理扫描二维码事件
type ScannerEventHandle struct {
}

func (this *ScannerEventHandle) GetKey() string {
	return "Scanner"
}
func (this *ScannerEventHandle) Handle(data []byte) interface{} {
	event := new(model.WeixinScannerEvent)
	e := xml.Unmarshal(data, event)
	if e != nil {
		log.Warnln(e)
		return nil
	}
	if weixinEventHanldeImpl != nil {
		return weixinEventHanldeImpl.OnScanner(event)
	}
	return nil
}

//处理上报地理位置事件
type LocationEventHandle struct {
}

func (this *LocationEventHandle) GetKey() string {
	return "location"
}
func (this *LocationEventHandle) Handle(data []byte) interface{} {
	event := new(model.WeixinLocationEvent)
	e := xml.Unmarshal(data, event)
	if e != nil {
		log.Warnln(e)
		return nil
	}
	if weixinEventHanldeImpl != nil {
		return weixinEventHanldeImpl.OnLocationEvent(event)
	}
	return nil
}

//处理自定义菜单事件
type MenuEventHandle struct {
}

func (this *MenuEventHandle) GetKey() string {
	return "menu"
}
func (this *MenuEventHandle) Handle(data []byte) interface{} {
	event := new(model.WeixinMenuEvent)
	e := xml.Unmarshal(data, event)
	if e != nil {
		log.Warnln(e)
		return nil
	}
	if weixinEventHanldeImpl != nil {
		return weixinEventHanldeImpl.OnMenu(event)
	}
	return nil
}
