package main

import (
	"go-dog/pkg/service"
	"go-dog/log"
	"go-dog/plugins"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/mem"
)

//GetCPUPercent 获取cpu
func GetCPUPercent() float64 {
	percent, _ := cpu.Percent(time.Second, false)
	return percent[0]
}

//GetMemPercent 获取内存
func GetMemPercent() float64 {
	memInfo, _ := mem.VirtualMemory()
	return memInfo.UsedPercent
}

//GetDiskPercent 获取磁盘
func GetDiskPercent() float64 {
	parts, _ := disk.Partitions(true)
	diskInfo, _ := disk.Usage(parts[0].Mountpoint)
	return diskInfo.UsedPercent
}

//EchoService 声明接口类
type EchoService struct {
	count int32
}

//RequestCount 请求数量统计
func (e *EchoService) RequestCount() {
	for {
		select {
		case <-time.After(time.Second * 5):
			var m runtime.MemStats
			runtime.ReadMemStats(&m)
			//log.Traceln("收到客户端请求:", atomic.LoadInt32(&e.count))
		}
	}
}

//Echo 定义方法Echo
func (e *EchoService) Echo(ctx plugins.Context, arg string) (string, error) {
	atomic.AddInt32(&e.count, 1)
	return "hello client", nil
}

//Arg 请求
type Arg struct {
	A int
	B int
}

//Add 定义方法Add
func (e *EchoService) Add(ctx plugins.Context, arg Arg) (int, error) {
	atomic.AddInt32(&e.count, 1)
	time.Sleep(1 * time.Second)
	return arg.A + arg.B, nil
}

//Set 请求
type Set struct {
	Value string
}

//Back 返回
type Back struct {
	Value string
}

//Set 定义方法Set
func (e *EchoService) Set(ctx plugins.Context, set Set) (Back, error) {
	atomic.AddInt32(&e.count, 1)
	return Back{
		Value: set.Value,
	}, nil
}

//IsOk 定义方法IsOk
func (e *EchoService) IsOk(ctx plugins.Context, ok bool) (bool, error) {
	atomic.AddInt32(&e.count, 1)
	return ok, nil
}

func main() {
	ser := service.CreateService(10)
	service := &EchoService{count: 0}
	ser.RegisterAPI("Add", "v1", "add", plugins.POST, 3, false, "测试Add", service.Add)
	ser.RegisterRPC("Echo", 3, false, "测试Echo", service.Echo)
	ser.RegisterRPC("Set", 3, false, "测试Set", service.Set)
	ser.RegisterRPC("IsOk", 3, false, "测试IsOk", service.IsOk)
	go service.RequestCount()
	err := ser.Run()
	if err != nil {
		log.Traceln("服务退出", err.Error())
	}
}
