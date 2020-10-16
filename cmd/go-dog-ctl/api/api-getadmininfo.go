package api

import (
	"encoding/json"
	"go-dog/cmd/define"
	"go-dog/cmd/go-dog-ctl/param"
	customerror "go-dog/error"
	"go-dog/plugins"
)

//GetAdminInfo 获取管理员信息
func (pointer *API) GetAdminInfo(ctx plugins.Context, request param.GetAdminInfoReq) (response param.GetAdminInfoRes, err error) {
	if e := json.Unmarshal([]byte(userInfo), &response); e != nil {
		err = customerror.EnCodeError(define.AdminLoginErr, e.Error())
		return
	}
	if e := json.Unmarshal([]byte(roleObj), &response.Role); e != nil {
		err = customerror.EnCodeError(define.AdminLoginErr, e.Error())
		return
	}
	return
}

const userInfo = `{
    "id": 1211313131231321,
    "name": "天野远子",
    "username": "admin",
    "password": "",
    "avatar": "/avatar2.jpg",
    "status": 1,
    "telephone": "",
    "lastLoginIp": "27.154.74.117",
    "lastLoginTime": 1534837621348,
    "creatorId": "admin",
    "createTime": 1497160610259,
    "merchantCode": "TLif2btpzg079h15bk",
    "deleted": 0,
    "roleId": "admin"
  }`

const roleObj = `{
    "id": "admin",
    "name": "管理员",
    "permissions": [{
      "roleId": "admin",
      "permissionId": "admin",
      "permissionName": "超级管理员",
      "actionEntitySet": [{
        "action": "add",
        "describe": "新增",
        "defaultCheck": false
      }, {
        "action": "query",
        "describe": "查询",
        "defaultCheck": false
      }, {
        "action": "get",
        "describe": "详情",
        "defaultCheck": false
      }, {
        "action": "update",
        "describe": "修改",
        "defaultCheck": false
      }, {
        "action": "delete",
        "describe": "删除",
        "defaultCheck": false
      }]
	}]
  }`
