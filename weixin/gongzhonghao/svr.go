package gongzhonghao

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"

	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/weixin/gongzhonghao/handle"
	wxmodel "github.com/tang-go/go-dog/weixin/gongzhonghao/model"
)

var gSess *WeiXinSession

func index(w http.ResponseWriter, r *http.Request) {
	var s []string
	var str string
	log.Debugln("------%v", r.URL.String())
	token := gSess.cfg.Token
	r.ParseForm() //解析参数, 默认是不会解析的
	s = append(s, token)
	signature := r.URL.Query().Get("signature")
	s = append(s, r.URL.Query().Get("timestamp"))
	s = append(s, r.URL.Query().Get("nonce"))
	echostr := r.URL.Query().Get("echostr")
	sort.Strings(s)
	for _, v := range s {
		str += v
	}
	//产生一个散列值得方式是 sha1.New()，sha1.Write(bytes)，然后 sha1.Sum([]byte{})。这里我们从一个新的散列开始。
	h := sha1.New()
	//写入要处理的字节。如果是一个字符串，需要使用[]byte(s) 来强制转换成字节数组。
	h.Write([]byte(str))
	//这个用来得到最终的散列值的字符切片。Sum 的参数可以用来都现有的字符切片追加额外的字节切片：一般不需要要。
	bs := h.Sum(nil)
	hex.EncodeToString(bs)

	//校验是否是微信服务器发来的消息
	if strings.Compare(signature, hex.EncodeToString(bs)) == 0 {
		//Get代表服务器做校验
		if strings.Compare(strings.ToLower(r.Method), "get") == 0 && (len(echostr) > 0) {
			w.Write([]byte(echostr))
		} else if strings.Compare(strings.ToLower(r.Method), "post") == 0 {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Errorln(err)
			}
			log.Debugln("\n微信消息:\n", string(body))
			//解析头部
			var msg wxmodel.WeixinBase
			if e := xml.Unmarshal(body, &msg); e != nil {
				log.Warnln(e)
			}
			rsp := handle.Entrance(msg, body)

			r, e := xml.Marshal(rsp)
			if e != nil {
				log.Warnln(string(r), " \n err:", e)
			}
			log.Debugln("\n返回微信消息:\n", string(r))
			w.Write(r)

		} else {
			w.Write(nil)
		}
	}

}

func svrInit(sess *WeiXinSession, cb handle.WeixinMsgHandleInterface, event handle.WeixinEventHandleInterface) {
	gSess = sess
	handle.Init(cb, event)
	http.HandleFunc("/", index) //设置访问的路由
	//后台启动
	go func() {
		log.Infoln("微信启动在", sess.cfg.Addr)
		err := http.ListenAndServe(sess.cfg.Addr, nil) //设置监听的端口
		if err != nil {
			log.Errorln("ListenAndServe: ", err)
		}
	}()

}
