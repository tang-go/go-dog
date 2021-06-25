module github.com/tang-go/go-dog

go 1.16

require (
	github.com/HdrHistogram/hdrhistogram-go v1.1.0 // indirect
	github.com/coreos/bbolt v1.3.4 // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf // indirect
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f // indirect
	github.com/dustin/go-humanize v1.0.0 // indirect
	github.com/gin-contrib/pprof v1.3.0 
	github.com/gin-gonic/gin v1.6.3
	github.com/go-redis/redis v6.15.9+incompatible
	github.com/go-sql-driver/mysql v1.5.0
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/golang/groupcache v0.0.0-20200121045136-8c9f03a8e57e // indirect
	github.com/google/btree v1.0.1 // indirect
	github.com/google/uuid v1.2.0 // indirect
	github.com/grpc-ecosystem/go-grpc-middleware v1.2.2 // indirect
	github.com/grpc-ecosystem/grpc-gateway v1.16.0 // indirect
	github.com/hashicorp/consul/api v1.8.1
	github.com/jinzhu/gorm v1.9.16
	github.com/json-iterator/go v1.1.10
	github.com/konsorten/go-windows-terminal-sequences v1.0.3
	github.com/nacos-group/nacos-sdk-go v1.0.7
	github.com/nsqio/go-nsq v1.0.8
	github.com/opentracing/opentracing-go v1.2.0
	github.com/prometheus/client_golang v1.10.0 
	github.com/satori/go.uuid v1.2.0
	github.com/sipt/GoJsoner v0.0.0-20170413020122-3e1341522aa6
	github.com/soheilhy/cmux v0.1.5 // indirect
	github.com/swaggo/gin-swagger v1.3.0
	github.com/tmc/grpc-websocket-proxy v0.0.0-20201229170055-e5319fda7802 // indirect
	github.com/uber/jaeger-client-go v2.25.0+incompatible
	github.com/uber/jaeger-lib v2.4.0+incompatible // indirect
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	go.uber.org/atomic v1.7.0 // indirect
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/sys v0.0.0-20210309074719-68d13333faf2
	golang.org/x/time v0.0.0-20201208040808-7e3f01d25324
	google.golang.org/genproto v0.0.0-20210113195801-ae06605f4595 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	sigs.k8s.io/yaml v1.2.0 // indirect
	xorm.io/core v0.7.3
	xorm.io/xorm v1.0.5
)

replace (
	github.com/coreos/bbolt v1.3.4 => go.etcd.io/bbolt v1.3.4
	github.com/mitchellh/cli v1.1.0 => github.com/mitchellh/cli v1.1.2
	google.golang.org/grpc => google.golang.org/grpc v1.26.0
)
