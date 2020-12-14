package utils

import (
	"encoding/base64"
	"strconv"

	"github.com/tang-go/go-dog/log"
)

const (
	StdPadding  rune = '-' // Standard padding character
	NoPadding   rune = -1  // No padding
	base64Table      = "123456789==3ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
)

var coder = base64.NewEncoding(base64Table).WithPadding(NoPadding)

func Base64Encode(src []byte) string {
	return coder.EncodeToString(src)
}

func Base64Decode(src []byte) (o string, err error) {
	str, err := coder.DecodeString(string(src))
	if err != nil {
		return "", err
	}
	o = string(str)
	return
}

func Base64EncodeIntToString(i int) string {
	str := strconv.Itoa(i)
	return coder.EncodeToString([]byte(str))
}

func Base64EncodeToString(i string) string {
	return coder.EncodeToString([]byte(i))
}

func Base64DecodeToInt(i string) (o int, err error) {
	src := []byte(i)
	str, err := coder.DecodeString(string(src))
	if err != nil {
		log.Errorln(err)
		return 0, err
	}

	temp := string(str)

	o, err = strconv.Atoi(temp)
	if err != nil {
		log.Errorln(err)
		return 0, err
	}

	return
}

func IsInt(i string) bool {
	_, o := strconv.Atoi(i)
	if o != nil {
		return false
	}
	return true
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue
		}
		ret = append(ret, a[i])
	}
	return
}

//debyte := base64Encode([]byte(hello))
//enbyte, err := base64Decode(debyte)
