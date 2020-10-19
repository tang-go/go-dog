package api

import (
	"go-dog/cmd/define"
	"go-dog/cmd/go-dog-ctl/param"
	"go-dog/cmd/go-dog-ctl/table"
	customerror "go-dog/error"
	"go-dog/plugins"
)

//GetAdminInfo 获取管理员信息
func (pointer *API) GetAdminInfo(ctx plugins.Context, request param.GetAdminInfoReq) (response param.GetAdminInfoRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	response.ID = admin.AdminID
	response.Name = admin.Name
	response.Avatar = "/avatar2.jpg"
	response.Phone = admin.Phone
	role := new(table.OwnerRole)
	if pointer.mysql.GetReadEngine().Where("role_id = ?", admin.RoleID).First(role).RecordNotFound() == true {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员权限不正确")
		return
	}
	response.RoleID = role.Name
	response.Role.ID = role.Name
	if role.IsAdmin {
		response.Role.Permissions = append(response.Role.Permissions, &param.Permissions{
			RoleID:         role.Name,
			PermissionID:   "admin",
			PermissionName: role.Description,
			ActionEntitySet: []*param.ActionEntitySet{
				&param.ActionEntitySet{
					Action:       "add",
					Describe:     "新增",
					DefaultCheck: true,
				},
				&param.ActionEntitySet{
					Action:       "delete",
					Describe:     "删除",
					DefaultCheck: true,
				},
				&param.ActionEntitySet{
					Action:       "update",
					Describe:     "修改",
					DefaultCheck: true,
				},
				&param.ActionEntitySet{
					Action:       "query",
					Describe:     "查询",
					DefaultCheck: true,
				},
			},
		})
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
