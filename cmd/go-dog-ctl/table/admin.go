package table

const (
	//LowAdmin 低级管理员
	LowAdmin = 0
	//GeneralAdmin 普通管理员---能添加管理
	GeneralAdmin = 1
	//SuperAdmin 超级管理员---拥有所有权限
	SuperAdmin = 2
)

//Admin 超级管理员
type Admin struct {
	//自增ID
	ID int64
	//账号自动生成--规程ID加随机6位数
	AdminID int64
	//名称
	Name string
	//电话
	Phone string
	//密码
	Pwd string
	//盐值 md5使用
	Salt string
	//等级
	Level int
	//所属业主
	OwnerID int64
	//是否被禁用
	IsDisable bool
	//注册时间
	Time int64
}
