module go-dog

go 1.13

require (
	github.com/StackExchange/wmi v0.0.0-20190523213315-cbe66965904d // indirect
	github.com/coreos/etcd v3.3.25+incompatible
	github.com/gin-gonic/gin v1.6.3
	github.com/go-ole/go-ole v1.2.4 // indirect
	github.com/micro/go-micro v1.18.0
	github.com/satori/go.uuid v1.2.0
	github.com/shirou/gopsutil v2.20.8+incompatible
	github.com/sipt/GoJsoner v0.0.0-20170413020122-3e1341522aa6
	github.com/sirupsen/logrus v1.7.0
	github.com/vmihailenco/msgpack v4.0.4+incompatible
	k8s.io/apimachinery v0.18.0
	k8s.io/client-go v0.18.0
	k8s.io/utils v0.0.0-20200912215256-4140de9c8800 // indirect
)

replace google.golang.org/grpc v1.27.0 => google.golang.org/grpc v1.26.0

replace github.com/sirupsen/logrus v1.7.0 => ./lib/sirupsen/logrus
