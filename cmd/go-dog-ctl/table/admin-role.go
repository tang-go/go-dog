package table

//AdminRole 管理员角色关联表
type AdminRole struct {
	//角色ID
	RoleID int64
	//管理员ID
	AdminID int64
	//创建时间
	Time int64
}
