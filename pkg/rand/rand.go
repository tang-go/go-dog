package rand

//随机数相关算法
import (
	"math/rand"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var _Chars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789")

//StringRand 生成随机字符串
func StringRand(length int) string {
	if length == 0 {
		return ""
	}
	clen := len(_Chars)
	if clen < 2 || clen > 256 {
		panic("Wrong charset length for NewLenChars()")
	}
	maxrb := 255 - (256 % clen)
	b := make([]byte, length)
	r := make([]byte, length+(length/4)) // storage for random bytes.
	i := 0
	for {
		if _, err := rand.Read(r); err != nil {
			panic("Error reading random bytes: " + err.Error())
		}
		for _, rb := range r {
			c := int(rb)
			if c > maxrb {
				continue
			}
			b[i] = _Chars[c%clen]
			i++
			if i == length {
				return string(b)
			}
		}
	}
}

//IntRand rand一个范围值
func IntRand(min, max int) int {
	if min >= max {
		return max
	}
	return rand.Intn(max-min) + min
}
