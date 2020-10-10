package md5

import (
	"crypto/md5"
	"fmt"
	"io"
)

//Md5 Md5加密
func Md5(in interface{}) string {
	h := md5.New()
	io.WriteString(h, fmt.Sprintf("%v", in))
	return fmt.Sprintf("%x", h.Sum(nil))
}
