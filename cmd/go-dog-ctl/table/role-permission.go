package table

//RolePermission 角色权限关联表
type RolePermission struct {
	//角色ID
	RoleID int64
	//权限ID
	PermissionID int64
	//时间
	Time int64
}
