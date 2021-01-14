module github.com/tang-go/go-dog

go 1.13

require (
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.10
	github.com/konsorten/go-windows-terminal-sequences v1.0.3
	github.com/nsqio/go-nsq v1.0.8
	github.com/opentracing/opentracing-go v1.2.0
	github.com/satori/go.uuid v1.2.0
	github.com/sipt/GoJsoner v0.0.0-20170413020122-3e1341522aa6
	github.com/stretchr/testify v1.6.1
	github.com/swaggo/gin-swagger v1.3.0
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/sys v0.0.0-20201214210602-f9fddec55a1e
	xorm.io/core v0.6.3
	xorm.io/xorm v1.0.5 // indirect
)

replace google.golang.org/grpc v1.27.0 => google.golang.org/grpc v1.26.0
