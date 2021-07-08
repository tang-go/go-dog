package nsq

import (
	"strings"
	"sync"
	"time"

	"github.com/nsqio/go-nsq"
	"github.com/tang-go/go-dog/lib/rand"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
)

//Logger 日志
type Logger struct {
}

//Output 输出
func (l *Logger) Output(calldepth int, s string) error {
	log.Traceln(s)
	return nil
}

//Nsq 生成者
type Nsq struct {
	cfg      plugins.Cfg
	pos      int
	count    int
	address  []string
	producer *nsq.Producer
	consumer []*nsq.Consumer
	lock     sync.Mutex
}

//NewNsq 新建一个Nsq
func NewNsq(cfg plugins.Cfg) *Nsq {
	n := new(Nsq)
	n.cfg = cfg
	n.pos = 0
	n.count = len(cfg.GetNsq())
	n.address = cfg.GetNsq()
	index := rand.IntRand(0, n.count)
	addr := n.address[index]
	if err := n.initProducer(addr); err != nil {
		log.Fatalf("连接%v失败: %s", n.cfg.GetNsq(), err.Error())
	}
	return n
}

func (n *Nsq) initProducer(addr string) error {
	n.lock.Lock()
	defer n.lock.Unlock()
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(addr, config)
	if err != nil {
		return err
	}
	producer.SetLogger(new(Logger), nsq.LogLevelInfo)
	err = producer.Ping()
	if err != nil {
		producer.Stop()
		return err
	}
	n.producer = producer
	return nil
}

func (n *Nsq) setConsumer(consumer *nsq.Consumer) {
	n.lock.Lock()
	n.consumer = append(n.consumer, consumer)
	n.lock.Unlock()
}

func (n *Nsq) getProducer() *nsq.Producer {
	n.lock.Lock()
	defer n.lock.Unlock()
	return n.producer
}

//Publish 发送消息
func (n *Nsq) Publish(topic string, msg []byte) error {
	if err := n.getProducer().Ping(); err != nil {
		n.getProducer().Stop()
		for _, addr := range n.address {
			//此处全部断线重连一下
			if err := n.initProducer(addr); err == nil {
				break
			}
		}
	}
	return n.getProducer().Publish(topic, msg)
}

//DeferredPublish 延迟消息
func (n *Nsq) DeferredPublish(topic string, delay time.Duration, msg []byte) error {
	if err := n.getProducer().Ping(); err != nil {
		n.getProducer().Stop()
		for _, addr := range n.address {
			//此处全部断线重连一下
			if err := n.initProducer(addr); err == nil {
				break
			}
		}
	}
	return n.getProducer().DeferredPublish(topic, delay, msg)
}

//Consumer 创建消费者
func (n *Nsq) Consumer(topic, channel string, f func(msg []byte) error) error {
	config := nsq.NewConfig()
	config.MaxAttempts = 65534
	config.LookupdPollInterval = time.Second
	config.MaxInFlight = len(n.cfg.GetNsq())
	topic = strings.Replace(topic, "/", "_", -1)
	channel = strings.Replace(channel, "/", "_", -1)
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		log.Errorf("创建消费者失败: %s", err.Error())
		return err
	}
	consumer.SetLogger(new(Logger), nsq.LogLevelInfo)
	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		if err := f(message.Body); err != nil {
			if !message.IsAutoResponseDisabled() {
				message.RequeueWithoutBackoff(time.Second * 5)
			}
		} else {
			message.Finish()
		}
		return nil
	}))
	//建立多个nsqd连接
	if err := consumer.ConnectToNSQDs(n.cfg.GetNsq()); err != nil {
		log.Errorf("连接%v失败: %s", n.cfg.GetNsq(), err.Error())
		return err
	}
	n.setConsumer(consumer)
	<-consumer.StopChan
	stats := consumer.Stats()
	log.Tracef("message received %d, finished %d, requeued:%s, connections:%s", stats.MessagesReceived, stats.MessagesFinished, stats.MessagesRequeued, stats.Connections)
	return nil
}

//Close 关闭
func (n *Nsq) Close() {
	for _, consumer := range n.consumer {
		consumer.Stop()
	}
	n.producer.Stop()
}
