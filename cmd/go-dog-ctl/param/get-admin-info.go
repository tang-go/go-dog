package param

//GetAdminInfoReq 获取管理员信息
type GetAdminInfoReq struct {
	Code string `json:"code" description:"业务随机码" type:"string"`
}

//GetAdminInfoRes 获取用户信息返回
type GetAdminInfoRes struct {
	ID     int64  `json:"id" description:"用户ID" type:"int64"`
	Name   string `json:"name" description:"名称" type:"string"`
	Avatar string `json:"avatar" description:"头像" type:"string"`
	Phone  string `json:"phone" description:"电话" type:"string"`
	RoleID string `json:"roleId" description:"权限名称" type:"string"`
	Role   Role   `json:"role" description:"权限内容" type:"Role"`
}

//Role 权限配置
type Role struct {
	ID          string         `json:"id" description:"权限名称" type:"string"`
	Permissions []*Permissions `json:"permissions" json:"id" description:"权限内容" type:"[]*Permissions "`
}

//Permissions 权限
type Permissions struct {
	RoleID          string             `json:"roleId" description:"权限名称" type:"string"`
	PermissionID    string             `json:"permissionId" description:"操作名称" type:"string"`
	PermissionName  string             `json:"permissionName" description:"操作描述" type:"string"`
	ActionEntitySet []*ActionEntitySet `json:"actionEntitySet" description:"操作动作" type:"[]*ActionEntitySet "`
}

//ActionEntitySet 动作
type ActionEntitySet struct {
	Action       string `json:"action" description:"动作名称" type:"string"`
	Describe     string `json:"describe" description:"动作描述" type:"string"`
	DefaultCheck bool   `json:"defaultCheck" description:"是否可以操作" type:"bool"`
}
