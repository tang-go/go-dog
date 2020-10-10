package table

//Menu 菜单
type Menu struct {
	//自增ID
	ID int64
	//菜单名称
	Title string
	//图标
	Icon string
	//路径
	Path string
	//等级---管理员可见等级
	Level int
	//母节点ID
	ParentID int64
	//时间
	Time int64
}
