package utils

import (
	"bytes"
	"errors"
	"math/rand"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/tang-go/go-dog/lib/encrypt"
)

const (
	SHORTTIMEFORMAT = "2006-01-02"
	LONGTIMEFORMAT  = "2006-01-02 15:04:05"
	MIDTIMEFORMAT   = "2006-01-02 15:04"
	TIGHTFORMAT     = "20060102150405"
	LOCALTIMEZONE   = "Asia/Shanghai"
)

func RandomString(l int) string {
	var result bytes.Buffer
	var temp string
	for i := 0; i < l; {
		temp = string(rune(65 + rand.Intn(25)))
		result.WriteString(temp)
		i++

	}
	return result.String()
}

//获取协程号
func GetGID() uint64 {
	b := make([]byte, 64)
	b = b[:runtime.Stack(b, false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, _ := strconv.ParseUint(string(b), 10, 64)
	return n
}

//合并map
func MergeMap(origin map[string]interface{}, add map[string]interface{}) {
	for k, v := range add {
		origin[k] = v
	}
}

//获取map key
func GetMapKeys(i map[string]interface{}) (o []string) {
	for k, _ := range i {
		o = append(o, k)
	}
	return
}

func GetMapKeysOfBoolValue(i map[string]bool) (o []string) {
	for k, _ := range i {
		o = append(o, k)
	}
	return
}

//includes
func Includes(key string, i []string) bool {
	for _, v := range i {
		if v == key {
			return true
		}
	}
	return false
}

//IncludesInt
func IncludesInt(key int, i []int) bool {
	for _, v := range i {
		if ToString(v) == ToString(key) {
			return true
		}
	}
	return false
}

//NotIncludes
func NotIncludes(key string, i []string) bool {
	for _, v := range i {
		if v == key {
			return false
		}
	}
	return true
}

//TimeShortFormat
func TimeShortFormat(i string) (time.Time, error) {
	loc, _ := time.LoadLocation(LOCALTIMEZONE)
	return time.ParseInLocation(SHORTTIMEFORMAT, i, loc)
}

//TimeLongFormat
func TimeLongFormat(i string) (time.Time, error) {
	loc, _ := time.LoadLocation(LOCALTIMEZONE)
	return time.ParseInLocation(LONGTIMEFORMAT, i, loc)
}

//ToString
func ToString(i int) string {
	return strconv.Itoa(i)
}

//ToInt
func ToInt(i string) (int, error) {
	return strconv.Atoi(i)
}

func GenOutTradeNum() (o string) {
	//length = 19
	return time.Now().Format(TIGHTFORMAT) + encrypt.GetRandomString(5)
}

//生成游客模式假电话
func GenVisitorFadePhone() string {
	return "v" + encrypt.GetRandomString(1) + time.Now().Format(TIGHTFORMAT)[5:14]
}

//生成client 的随机openid
func GenClientFadeOpenid() string {
	return "c" + encrypt.GetRandomString(1) + time.Now().Format(TIGHTFORMAT)[5:14]
}
func TimeAdd(i int, t string) time.Time {
	now := time.Now()

	switch t {
	case "d":
		return now.AddDate(0, 0, i)
	case "h":
		num := i
		str := strconv.Itoa(num) + t
		mm, _ := time.ParseDuration(str)
		return now.Add(mm)
	case "m":
		num := i
		str := strconv.Itoa(num) + t
		mm, _ := time.ParseDuration(str)
		return now.Add(mm)
	}
	return now
}

func OmitMapEmptyKey(o map[string]interface{}) {
	for k, v := range o {
		if v == "" {
			delete(o, k)
		}
	}
}

//删除slice元素

func DeleteSlice(slice interface{}, index int) (interface{}, error) {
	sliceValue := reflect.ValueOf(slice)
	length := sliceValue.Len()
	if slice == nil || length == 0 || (length-1) < index {
		return nil, errors.New("error")
	}
	if length-1 == index {
		return sliceValue.Slice(0, index).Interface(), nil
	} else if (length - 1) >= index {
		return reflect.AppendSlice(sliceValue.Slice(0, index), sliceValue.Slice(index+1, length)).Interface(), nil
	}
	return nil, errors.New("error")
}

//数组去重
func RemoveDupByMap(i []string) []string {
	o := make([]string, 0)
	tempMap := make(map[string]byte) // 存放不重复主键
	for _, e := range i {
		if _, has := tempMap[e]; !has {
			tempMap[e] = 1
			o = append(o, e)
		}
	}
	return o
}

const SQL_FORMAT = "2006-01-02 15:04:05 "

func DayBefore(days int) (start time.Time, end time.Time) {
	day := time.Now().Add(-time.Hour * time.Duration(24*days))
	start = DayBegin(day)
	end = DayEnd(day)
	return start, end
}

func DayBegin(t time.Time) time.Time {
	return t.Add(-time.Duration(t.Hour()) * time.Hour).Add(-time.Duration(t.Minute()) * time.Minute).Add(-time.Duration(t.Second()) * time.Second)
}

func DayEnd(t time.Time) time.Time {
	return t.Add(time.Duration(24-t.Hour()) * time.Hour).Add(-time.Duration(t.Minute()) * time.Minute).Add(-time.Duration(t.Second()) * time.Second)
}
