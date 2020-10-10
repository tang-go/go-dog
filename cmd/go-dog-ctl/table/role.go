package table

const (
	//IsAdmin 是超级管理员
	IsAdmin = true
	//NoAdmin 不是管理员
	NoAdmin = false
)

//OwnerRole 业主角色表
type OwnerRole struct {
	//角色ID
	RoleID int64
	//角色名称
	Name string
	//角色描述
	Description string
	//是否为超级管理员
	IsAdmin bool
	//业主ID
	OwnerID int64
	//角色创建时间
	Time int64
}
