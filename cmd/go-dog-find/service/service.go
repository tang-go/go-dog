package service

import (
	"fmt"
	"go-dog/cache"
	"go-dog/cmd/go-dog-find/param"
	"go-dog/internal/config"
	"go-dog/lib/io"
	"go-dog/lib/uuid"
	"go-dog/log"
	"go-dog/plugins"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

//Service 控制服务
type Service struct {
	cfg        plugins.Cfg
	registers  map[string]*Register
	discoverys map[string]*Discovery
	cache      *cache.Cache
	close      int32
	lock       sync.RWMutex
}

//NewService 初始化服务
func NewService() *Service {
	s := &Service{
		close:      0,
		registers:  make(map[string]*Register),
		discoverys: make(map[string]*Discovery),
	}
	s.cfg = config.NewConfig()
	//初始化缓存
	s.cache = cache.NewCache(s.cfg)
	//初始化日志
	switch s.cfg.GetRunmode() {
	case "panic":
		log.SetLevel(log.PanicLevel)
		break
	case "fatal":
		log.SetLevel(log.FatalLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "warn":
		log.SetLevel(log.WarnLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "trace":
		log.SetLevel(log.TraceLevel)
		break
	default:
		log.SetLevel(log.TraceLevel)
		break
	}
	return s
}

//Run 启动
func (s *Service) Run() error {
	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		err := s._Run()
		if err != nil {
			panic(err.Error())
		}
	}()
	msg := <-c
	s.Close()
	return fmt.Errorf("收到kill信号:%s", msg)
}

//Close 关闭服务
func (s *Service) Close() {
	atomic.AddInt32(&s.close, 1)
}

//_Run
func (s *Service) _Run() error {
	address := fmt.Sprintf("0.0.0.0:%d", s.cfg.GetPort())
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	log.Traceln("服务发现中心启动:", address)
	defer l.Close()
	for {
		if atomic.LoadInt32(&s.close) > 0 {
			return nil
		}
		conn, err := l.Accept()
		if err != nil {
			log.Errorln(err.Error())
			continue
		}
		go s._ServeConn(conn)
	}
}

// ServeConn 拦截一个链接
func (s *Service) _ServeConn(conn net.Conn) {
	_, buff, err := io.ReadByTime(conn, time.Now().Add(time.Second*5))
	if err != nil {
		conn.Close()
		log.Errorln(err.Error())
		return
	}
	login := new(param.LoginReq)
	if err := login.DeCode(buff, login); err == nil {
		//服务发现
		if login.Type == param.DisType {
			token := uuid.GetToken()
			discovery := NewDiscovery(s, conn, func() {
				s.lock.Lock()
				delete(s.discoverys, token)
				s.lock.Unlock()
				log.Traceln("收到客户端下线:", token)
			})
			go discovery.Run()
			s.lock.Lock()
			s.discoverys[token] = discovery
			s.lock.Unlock()
			log.Traceln("收到客户端上线:", token)
		}
		//服务注册
		if login.Type == param.RegType {
			token := uuid.GetToken()
			register := NewRegister(s, conn, login.Data, func() {
				s.lock.Lock()
				delete(s.registers, token)
				s.lock.Unlock()
			})
			go register.Run()
			s.lock.Lock()
			s.registers[token] = register
			s.lock.Unlock()
		}
	}
}
