package table

//AdminMenu 管理员菜单
type AdminMenu struct {
	//自增ID
	ID int64
	//菜单ID
	MenuID int64
	//管理员ID--账号id
	AdminID int64
	//路由
	Path string
	//时间
	Time int64
}
