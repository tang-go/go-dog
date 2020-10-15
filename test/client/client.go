package main

import (
	"fmt"
	customerror "go-dog/error"
	"go-dog/internal/client"
	"go-dog/internal/config"
	"go-dog/internal/context"
	"go-dog/log"
	"go-dog/plugins"
	"sync"
	"sync/atomic"
	"time"
)

var gWait sync.WaitGroup

//Arg 请求
type Arg struct {
	A int
	B int
}

//Set 请求
type Set struct {
	Value string
}

//Back 返回
type Back struct {
	Value string
}

func main() {
	cli := client.NewClient(10, config.NewConfig())
A:
	now := time.Now()
	var count int32
	var errnum int32
	rquestnum := 1000
	for n := 0; n < rquestnum; n++ {
		gWait.Add(1)
		go func(id int) {
			defer gWait.Done()

			ctx := context.Background()
			ctx.SetAddress("127.0.0.1")
			ctx.SetTraceID(fmt.Sprintf("%d", id))
			ctx.SetIsTest(true)
			ctx = context.WithTimeout(ctx, int64(time.Second*4))
			var back bool
			arg := true
			err := cli.Call(ctx, plugins.RandomMode, "test", "IsOk", arg, &back)
			if err != nil {
				myError := customerror.DeCodeError(err)
				log.Errorln(myError.Code, myError.Msg)
				atomic.AddInt32(&errnum, 1)
			} else {
				//log.Traceln("收到返回结果:", back)
			}
			atomic.AddInt32(&count, 1)
		}(n)

	}
	gWait.Wait()
	log.Tracef("| 总数:%d | 失败:%d | 比例:%d | tps:%f |", count, errnum, errnum*100/count, float64(rquestnum)/time.Now().Sub(now).Seconds())
	time.Sleep(time.Second * 1)
	goto A
	//select {}
}
