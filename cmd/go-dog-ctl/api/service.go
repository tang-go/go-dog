package api

import (
	"fmt"
	"go-dog/cache"
	"go-dog/cmd/define"
	"go-dog/cmd/go-dog-ctl/table"
	"go-dog/internal/service"
	"go-dog/mysql"
	"go-dog/pkg/log"
	"go-dog/pkg/md5"
	"go-dog/pkg/rand"
	"go-dog/pkg/snowflake"
	"go-dog/plugins"
	"go-dog/serviceinfo"
	"math/big"
	"net"
	"reflect"
	"strconv"
	"sync"
	"time"
)

//APIService API服务
type _APIService struct {
	method  *serviceinfo.API
	name    string
	explain string
	count   int32
}

//Service 控制服务
type Service struct {
	service   plugins.Service
	mysql     *mysql.Mysql
	snowflake *snowflake.SnowFlake
	cache     *cache.Cache
	apis      map[string]*_APIService
	services  map[string]*serviceinfo.APIServiceInfo
	lock      sync.RWMutex
}

//NewService 初始化服务
func NewService() *Service {
	ctl := new(Service)
	ctl.apis = make(map[string]*_APIService)
	ctl.services = make(map[string]*serviceinfo.APIServiceInfo)
	//初始化rpc服务端
	ctl.service = service.CreateService(define.TTL)
	//设置服务名称
	ctl.service.SetName(define.SvcController)
	//设置服务端最大访问量
	ctl.service.GetLimit().SetLimit(define.MaxServiceRequestCount)
	//设置客户端最大访问量
	ctl.service.GetClient().GetLimit().SetLimit(define.MaxClientRequestCount)
	//注册API上线通知
	ctl.service.GetClient().GetDiscovery().RegAPIServiceOnlineNotice(ctl._ApiServiceOnline)
	//注册API下线通知
	ctl.service.GetClient().GetDiscovery().RegAPIServiceOfflineNotice(ctl._ApiServiceOffline)
	//开始监听API事件
	ctl.service.GetClient().GetDiscovery().WatchAPIService()
	//初始化数据库
	ctl.mysql = mysql.NewMysql(ctl.service.GetCfg())
	//初始化数据库表
	ctl.mysql.GetWriteEngine().AutoMigrate(
		table.Admin{},
		table.Owner{},
		table.OwnerRole{},
		table.Permission{},
		table.RolePermission{},
		table.AdminRole{},
		table.Log{},
	)
	//初始化缓存
	ctl.cache = cache.NewCache(ctl.service.GetCfg())
	//初始化雪花算法
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ctl.service.GetCfg().GetHost()).To4())
	id, err := strconv.ParseInt(fmt.Sprintf("%d%d", ret.Int64(), ctl.service.GetCfg().GetPort()), 10, 64)
	if err != nil {
		panic(err)
	}
	ctl.snowflake = snowflake.NewSnowFlake(id)
	//初始化API
	ctl.InitAPI()
	//初始化数据库数据
	ctl._InitMysql("13688460148", "admin")
	return ctl
}

//_Duplicate 去重
func _Duplicate(a interface{}) (ret []int64) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface().(int64))
	}
	return ret
}

//Run 启动
func (pointer *Service) Run() error {
	return pointer.service.Run()
}

//InitAPI 初始化API
func (pointer *Service) InitAPI() {
	//获取图片验证码
	pointer.service.GET("GetCode", "v1", "get/code",
		3,
		false,
		"获取图片验证码",
		pointer.GetCode)
	//验证码验证码
	pointer.service.POST("AdminLogin", "v1", "admin/login",
		3,
		false,
		"管理员登录",
		pointer.AdminLogin)
	//获取API列表
	pointer.service.GET("GetAPIList", "v1", "get/api/list",
		3,
		true,
		"获取api列表",
		pointer.GetAPIList)
	//获取服务列表
	pointer.service.GET("GetServiceList", "v1", "get/service/list",
		3,
		true,
		"获取服务列表",
		pointer.GetServiceList)
}

//_InitMysql 第一次加载初始化数据库数据
func (pointer *Service) _InitMysql(phone, pwd string) {
	//读取是否有业主了
	owner := new(table.Owner)
	if pointer.mysql.GetReadEngine().Where("phone = ?", phone).First(owner).RecordNotFound() == false {
		return
	}
	//如果没有业主则新增默认业主
	owner.OwnerID = pointer.snowflake.GetID()
	owner.Name = "超级业主"
	owner.Phone = phone
	owner.Level = 1
	owner.IsDisable = table.OwnerAvailable
	owner.IsAdminOwner = table.IsAdminOwner
	owner.Time = time.Now().Unix()
	//超级管理员
	ownerRole := &table.OwnerRole{
		RoleID: pointer.snowflake.GetID(),
		//角色名称
		Name: "超级管理员",
		//角色描述
		Description: "系统自带的超级管理员",
		//是否为超级管理员
		IsAdmin: table.IsAdmin,
		//业主ID
		OwnerID: owner.OwnerID,
		//角色创建时间
		Time: owner.Time,
	}
	//管理员
	admin := &table.Admin{
		//账号 唯一主键
		AdminID: pointer.snowflake.GetID(),
		//名称
		Name: "超级管理员",
		//电话
		Phone: phone,
		//盐值 md5使用
		Salt: rand.StringRand(6),
		//等级
		Level: owner.Level,
		//所属业主
		OwnerID: owner.OwnerID,
		//是否被禁用
		IsDisable: table.AdminAvailable,
		//注册事件
		Time: owner.Time,
	}
	//生成密码
	admin.Pwd = md5.Md5(md5.Md5(pwd) + admin.Salt)
	//生成权限映射
	adminRole := &table.AdminRole{
		//角色ID
		RoleID: ownerRole.RoleID,
		//管理员ID
		AdminID: admin.AdminID,
		//创建时间
		Time: owner.Time,
	}
	//开启数据库操作
	tx := pointer.mysql.GetWriteEngine().Begin()
	if err := tx.Create(owner).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Create(ownerRole).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Create(admin).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	if err := tx.Create(adminRole).Error; err != nil {
		tx.Rollback()
		panic(err)
	}
	tx.Commit()
}

//apiServiceOnline api服务上线
func (pointer *Service) _ApiServiceOnline(key string, service *serviceinfo.APIServiceInfo) {
	pointer.lock.Lock()
	for _, method := range service.API {
		url := "/api/" + service.Name + "/" + method.Version + "/" + method.Path
		if api, ok := pointer.apis[url]; ok {
			api.count++
		} else {
			pointer.apis[url] = &_APIService{
				method:  method,
				name:    service.Name,
				explain: service.Explain,
				count:   1,
			}
		}
		pointer.services[key] = service
	}
	pointer.lock.Unlock()
}

//apiServiceOffline api服务下线
func (pointer *Service) _ApiServiceOffline(key string) {
	pointer.lock.Lock()
	if service, ok := pointer.services[key]; ok {
		for _, method := range service.API {
			url := "/api/" + service.Name + "/" + method.Version + "/" + method.Path
			if api, ok := pointer.apis[url]; ok {
				api.count--
				if api.count <= 0 {
					delete(pointer.apis, url)
				}
			}
		}
		delete(pointer.services, key)
	}
	pointer.lock.Unlock()
}

// Set 设置验证码ID
func (pointer *Service) Set(id string, value string) {
	log.Traceln("Set", id)
	if err := pointer.cache.GetCache().SetByTime(id, value, define.CodeValidityTime); err != nil {
		log.Errorln(err.Error())
	}
}

// Get 更具验证ID获取验证码
func (pointer *Service) Get(id string, clear bool) (vali string) {
	log.Traceln("Get", id, clear)
	err := pointer.cache.GetCache().Get(id, &vali)
	if err != nil {
		log.Errorln(err.Error())
	}
	if clear {
		pointer.cache.GetCache().Del(id)
	}
	return
}

//Verify 验证验证码
func (pointer *Service) Verify(id, answer string, clear bool) bool {
	vali := pointer.Get(id, clear)
	if vali != answer {
		return false
	}
	return true
}
