package table

const (
	//AdminDisable  管理员禁用
	AdminDisable = true
	//AdminAvailable 管理员可用
	AdminAvailable = false
)

//Admin 超级管理员
type Admin struct {
	//账号 唯一主键
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
	//角色
	RoleID int64
	//注册事件
	Time int64
}
