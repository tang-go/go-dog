package table

const (
	//OwnerDisable  业主禁用
	OwnerDisable = true
	//OwnerAvailable 业主可用
	OwnerAvailable = false
	//IsAdminOwner 是超级业主
	IsAdminOwner = true
	//NoAdminOwner 不是超级业主
	NoAdminOwner = false
)

//Owner 业主
type Owner struct {
	//业主ID
	OwnerID int64
	//名称
	Name string
	//电话
	Phone string
	//Level 等级
	Level int
	//是否被禁用
	IsDisable bool
	//是否是超级业主
	IsAdminOwner bool
	//注册时间
	Time int64
}
