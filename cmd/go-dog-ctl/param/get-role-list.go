package param

//GetRoleListReq 获取角色列表请求
type GetRoleListReq struct {
	PageNo   int `json:"pageNo" description:"当前请求页数" type:"int"`
	PageSize int `json:"pageSize" description:"每一页的总数" type:"int"`
}

//GetRoleListRes 获取角色列表响应
type GetRoleListRes struct {
	Data       []RoleInfo `json:"data" description:"数据" type:"[]RoleInfo"`
	PageSize   int        `json:"pageSize" description:"一页大小" type:"int"`
	PageNo     int        `json:"pageNo" description:"当前所处的页数" type:"int"`
	TotalPage  int        `json:"totalPage" description:"总页数" type:"int"`
	TotalCount int        `json:"totalCount" description:"总数量" type:"int"`
}

//RoleInfo 角色信息
type RoleInfo struct {
	//角色ID
	RoleID int64 `json:"roleId" description:"角色ID" type:"int64"`
	//角色名称
	Name string `json:"name" description:"角色名称" type:"string"`
	//角色描述
	Description string `json:"description" description:"角色描述" type:"string"`
	//是否为超级管理员
	IsAdmin bool `json:"isAdmin" description:"是否为超级管理员" type:"bool"`
	//业主ID
	OwnerID int64 `json:"ownerId" description:"业主ID" type:"int64"`
	//角色创建时间
	Time string `json:"time" description:"创建时间" type:"string"`
}
