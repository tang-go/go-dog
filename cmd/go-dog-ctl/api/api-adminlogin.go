package api

import (
	"go-dog/cmd/define"
	"go-dog/cmd/go-dog-ctl/param"
	"go-dog/cmd/go-dog-ctl/table"
	customerror "go-dog/error"
	"go-dog/lib/md5"
	"go-dog/plugins"
	"time"
)

//AdminLogin 管理员登录
func (pointer *API) AdminLogin(ctx plugins.Context, request param.AdminLoginReq) (response param.AdminLoginRes, err error) {
	//查询是否拥有此用户
	admin := new(table.Admin)
	if pointer.mysql.GetReadEngine().Where("phone = ?", request.Phone).First(admin).RecordNotFound() == true {
		err = customerror.EnCodeError(define.AdminLoginErr, "管理员登录失败")
		return
	}
	//密码对比
	if md5.Md5(md5.Md5(request.Pwd)+admin.Salt) != admin.Pwd {
		err = customerror.EnCodeError(define.AdminLoginErr, "管理员登录失败")
		return
	}
	//生成登录记录
	mysqllog := &table.Log{
		//日志ID
		LogID: pointer.snowflake.GetID(),
		//类型
		Type: table.LoginType,
		//操作人
		AdminID: admin.AdminID,
		//名称
		AdminName: admin.Name,
		//操作方法
		Method: "AdminLogin",
		//描述
		Description: "管理员登录",
		//业主ID
		OwnerID: admin.OwnerID,
		//操作IP
		IP: ctx.GetAddress(),
		//操作URL
		URL: ctx.GetDataByKey("URL").(string),
		//操作时间
		Time: time.Now().Unix(),
	}
	if e := pointer.mysql.GetWriteEngine().Create(mysqllog).Error; e != nil {
		err = customerror.EnCodeError(define.AdminLoginErr, e.Error())
		return
	}
	//生成token
	token := md5.Md5(admin.AdminID)
	//生成token缓存
	pointer.cache.GetCache().SetByTime(token, admin, define.AdminTokenValidityTime)
	//登录成功返回
	response.Name = admin.Name
	response.OwnerID = admin.OwnerID
	response.Token = token
	return
}
