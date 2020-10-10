package param

//GetMenuRequest 获取菜单请求
type GetMenuRequest struct {
	//发起请求用户的token
	Token string
}

//GetMenuResponse 获取菜单响应
type GetMenuResponse struct {
	//菜单列表
	Menus []*Menu
}

//AddMenuRequest 添加菜单菜单请求
type AddMenuRequest struct {
	//发起请求用户的token
	Token string
	//路径
	Path string
	//名称
	Title string
	//图标
	Icon string
	//可见等级
	Level int
	//母节点ID
	ParentID int64
}

//AddMenuResponse 添加菜单菜单响应
type AddMenuResponse struct {
	//添加成功的menuID
	MenuID int64
}

//DelMenuRequest 删除菜单菜单请求
type DelMenuRequest struct {
	//发起请求用户的token
	Token string
	//母节点ID
	MenuID int64
}

//DelMenuResponse 删除菜单菜单响应
type DelMenuResponse struct {
	//成功true 失败false
	Suceess bool
}

//Menu 菜单
type Menu struct {
	//ID
	ID int64 `json:"id,omitempty"`
	//路径
	Path string `json:"path,omitempty"`
	//名称
	Title string `json:"title,omitempty"`
	//图标
	Icon string `json:"icon,omitempty"`
	//子菜单
	Children []*Menu `json:"children,omitempty"`
}
