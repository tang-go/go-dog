package net

import (
	"bytes"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"os"
	"time"
	"crypto/x509"
	"log"
)

type httpResponse struct {
	R    *http.Response
	Body []byte
}

func HttpsGet(url string) ([]byte, error) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal( err.Error(), url)
	}
	return b, err
}

func HttpGet(url string) (body []byte) {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return body
}

func HttpsPost(url ,key, cert , root string,data []byte)(body []byte){
	var tr *http.Transport
	certs, err := tls.LoadX509KeyPair(cert, key)
	if err != nil {
		log.Println("certs load err:", err)

	}
	rootCa, err := ioutil.ReadFile(root)
	if err != nil {
		log.Println("err2222:", err)
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(rootCa)
	tr = &http.Transport{
		TLSClientConfig: &tls.Config{
			RootCAs:      pool,
			Certificates: []tls.Certificate{certs},
		},
	}
	c := &http.Client{Transport: tr}
	resp, err := c.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(data))
	if err != nil {
		log.Fatal(err)
		return []byte("")
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
		return []byte("")
	}
	return body
}

func HttpPost(url string, data []byte) (body []byte) {
	client := http.Client{
		Transport:&http.Transport{
			TLSClientConfig:&tls.Config{
				InsecureSkipVerify:true,
			},
		},
	}
	resp, err := client.Post(url, "application/json;charset=utf-8", bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
	}
	return body
}

func HttpPostWithJson(url string, data string) (rsp *httpResponse, err error) {
	client := http.Client{}
	rsp = new(httpResponse)
	req_new := bytes.NewBuffer([]byte(data))
	request, _ := http.NewRequest("POST", url, req_new)
	request.Header.Set("Content-type", "application/json")
	rsp.R, err = client.Do(request)
	if rsp.R != nil {
		rsp.Body, err = ioutil.ReadAll(rsp.R.Body)
	}
	return
}

func Http(method string, url string) (body []byte, err error) {
	c := &http.Client{}

	//提交请求
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return body, nil
}

func HttpXml(method string, url string, data []byte) (body []byte, err error) {
	c := &http.Client{}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(data))
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	req.Header.Set("Accept", "application/xml")
	req.Header.Set("Content-Type", "application/xml;charset=utf-8")

	resp, err := c.Do(req)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return body, nil
}

func HttpWithBaseAuth(url, method, username, pass string) (body []byte, err error) {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.SetBasicAuth(username, pass)
	resp, err := client.Do(req)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	if bodyText, err := ioutil.ReadAll(resp.Body); err == nil {
		return bodyText, nil
	} else {
		//log.Fatal(err)
		return nil, err
	}
	return nil, nil
}

func HttpWithAuth(url, method, auth string) (body []byte, err error) {

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)
	req.Header.Set("Authorization", auth)
	req.Header.Set("Accept", "*/*")
	resp, err := client.Do(req)
	if err != nil {
		//log.Fatal(err)
		return nil, err
	}
	if bodyText, err := ioutil.ReadAll(resp.Body); err == nil {
		return bodyText, nil
	} else {
		//log.Fatal(err)
		return nil, err
	}
	return nil, nil
}

//使用本地文件上传
func UploadFile(uri string, params map[string]string, paramName, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, path)
	if err != nil {
		return err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return err
	}
	request, err := http.NewRequest("POST", uri, body)
	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return err
	}

	// Check the response
	if res.StatusCode != http.StatusOK {
		err = fmt.Errorf("bad status: %s", res.Status)
		return err
	}
	return err
}

const (
	ByteLenOf2M = 2097152 //2M字节大小
)

/**
*使用网络图片上传（将网络图片下载到内存然后转发）
* 文件名随机生成
*desUrl:上传路径
*sourceUrl：下载路径
*paramField：上传字段名
*paramName：上传字段名对应文件名
**/
func UploadPicByNetUrlSource(desUrl, sourceUrl, paramField string) ([]byte, error) {
	resp, err := http.Get(sourceUrl)
	defer resp.Body.Close()
	pix, err := ioutil.ReadAll(resp.Body)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	paramName := GetRandomString(10) + ".jpg"
	part, err := writer.CreateFormFile(paramField, paramName)
	if err != nil {
		return nil, err
	}
	byteLen, err := io.Copy(part, bytes.NewReader(pix))
	if err != nil {
		return nil, err
	}

	if byteLen > ByteLenOf2M {
		return nil, errors.New("文件超过2M，请使用小图重新编辑图文封面！")
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("POST", desUrl, body)
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()
	t, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return t, nil
}

//生成随机字符串
func GetRandomString(size int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
