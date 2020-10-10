package table

const (
	//LoginType 登录类型
	LoginType int32 = iota
)

//Log 日志
type Log struct {
	//日志ID
	LogID int64
	//类型
	Type int32
	//操作人
	AdminID int64
	//名称
	AdminName string
	//操作方法
	Method string
	//描述
	Description string
	//业主ID
	OwnerID int64
	//操作IP
	IP string
	//操作URL
	URL string
	//操作时间
	Time int64
}
