package gongzhonghao

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
)

const (
	tokenKey          = "access_token"
	ticketKey         = "ticket"
	expireIn          = 7200
	tokenURLFormatStr = "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s"
	signURLFormatStr  = "https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi"
)

// WxSigner is the struct for generating weixin signature
type WxSigner struct {
	TokenURL           string
	SignatureURLFormat string
	Ticket             string
	ExpireTime         int64
}

// Response contains values for configuring wx.config in front end
// see: http://mp.weixin.qq.com/wiki/7/aaa137b55fb2e0456bf8dd9148dd613f.html
type Response struct {
	NonceStr  string `json:"nonceStr"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

// before calling this method, NonceStr and Timestamp must be set
func (s *Response) genSignature(ticket, url string) {
	str := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s",
		ticket, s.NonceStr, s.Timestamp, url)
	h := sha1.New()
	h.Write([]byte(str))
	s.Signature = hex.EncodeToString(h.Sum(nil))
}

// NewSigner creates a WxSigner. AppID and secret are mandatory
func NewSigner(appID, secret string, urls ...string) *WxSigner {
	tokenURLFormat := tokenURLFormatStr
	signatureURLFormat := signURLFormatStr
	for i, v := range urls {
		if i == 0 {
			tokenURLFormat = v
		}
		if i == 1 {
			signatureURLFormat = v
		}
	}
	return &WxSigner{
		TokenURL:           fmt.Sprintf(tokenURLFormat, appID, secret),
		SignatureURLFormat: signatureURLFormat,
	}
}

// GenSignature generates the signature response for a given url
func (w *WxSigner) GenSignature(url string) (Response, error) {
	resp := Response{
		NonceStr:  randStr(16),
		Timestamp: time.Now().Unix(),
	}

	ticket, err := w.RequestTicket()
	if err != nil {
		return resp, err
	}

	fmt.Println("------> wx ticket", ticket)
	resp.genSignature(ticket, url)

	return resp, nil
}

// RequestTicket sends http requests to weixin server to fetch json api ticket
// In most cases, users don't need to call this method
func (w *WxSigner) RequestTicket() (string, error) {
	current := time.Now().Unix()
	// before the real expire time (exptime - margin), a new ticket will be requested
	// maybe it's not necessary
	if w.Ticket != "" && w.ExpireTime > current {
		return w.Ticket, nil
	}

	token, err := getResponse(w.TokenURL, tokenKey)
	if err != nil {
		return "", err
	}

	signatureURL := fmt.Sprintf(w.SignatureURLFormat, token)
	ticket, err := getResponse(signatureURL, ticketKey)
	if err != nil {
		return "", err
	}

	w.Ticket = ticket
	// 'current' is not precise, since the two requests will take some time
	// Then 'ExpireTime' will be slightly smaller than the real one
	// This situation is totally OK
	w.ExpireTime = current + expireIn

	return ticket, nil
}

func getResponse(url string, key string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result map[string]interface{}
	if err = json.Unmarshal(body, &result); err != nil {
		return "", err
	}

	valueInterface := result[key]
	value, ok := valueInterface.(string)
	if !ok || value == "" {
		return "", fmt.Errorf("fail to get correct response: url(%s) resp(%v)", url, result)
	}

	return value, nil
}

// borrow from stackoverflow
func randStr(n int) string {
	letterBytes := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	letterIdxBits := uint(6)                     // 6 bits to represent a letter index
	letterIdxMask := int64(1<<letterIdxBits - 1) // All 1-bits, as many as letterIdxBits
	letterIdxMax := 63 / letterIdxBits

	b := make([]byte, n)
	// A rand.Int63() generates 63 random bits, enough for letterIdxMax letters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return string(b)
}
