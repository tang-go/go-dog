package table

//Permission 管理员权限
type Permission struct {
	//权限ID
	PermissionID int64
	//角色名称
	Name string
	//角色描述
	Description string
	//权限接口链接
	URL string
	//业主ID
	OwnerID int64
	//角色创建时间
	Time int64
}
