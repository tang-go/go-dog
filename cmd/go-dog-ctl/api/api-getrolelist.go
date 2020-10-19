package api

import (
	"go-dog/cmd/define"
	"go-dog/cmd/go-dog-ctl/param"
	"go-dog/cmd/go-dog-ctl/table"
	customerror "go-dog/error"
	"go-dog/plugins"
	"time"
)

//GetRoleList 获取权限列表
func (pointer *API) GetRoleList(ctx plugins.Context, request param.GetRoleListReq) (response param.GetRoleListRes, err error) {
	admin, ok := ctx.GetShareByKey("Admin").(*table.Admin)
	if ok == false {
		err = customerror.EnCodeError(define.GetAdminInfoErr, "管理员信息失败")
		return
	}
	var roles []table.OwnerRole
	//err = pointer.GetReadEngine().Where("game_user_id = ? AND status != ?", gameUserID, model.Deleted).Order("status ASC").Order("created_at DESC").Find(&mails).Error
	if e := pointer.mysql.GetReadEngine().Where("owner_id = ?", admin.OwnerID).Limit(request.PageSize).Offset((request.PageNo - 1) * request.PageSize).Find(&roles).Error; e != nil {
		err = customerror.EnCodeError(define.GetRoleListErr, e.Error())
		return
	}
	total := 0
	if e := pointer.mysql.GetReadEngine().Model(&table.OwnerRole{}).Where("owner_id = ?", admin.OwnerID).Count(&total).Error; e != nil {
		err = customerror.EnCodeError(define.GetRoleListErr, e.Error())
		return
	}
	response.PageNo = request.PageNo
	response.PageSize = request.PageSize
	if total%request.PageSize > 0 {
		response.TotalPage = total/request.PageSize + 1
	}
	if total%request.PageSize < 0 {
		response.TotalPage = total / request.PageSize
	}
	response.TotalCount = total
	for _, role := range roles {
		response.Data = append(response.Data, param.RoleInfo{
			//角色ID
			RoleID: role.RoleID,
			//角色名称
			Name: role.Name,
			//角色描述
			Description: role.Description,
			//是否为超级管理员
			IsAdmin: role.IsAdmin,
			//业主ID
			OwnerID: role.OwnerID,
			//角色创建时间
			Time: time.Unix(role.Time, 0).Format("2006-01-02 15:04:05"),
		})
	}
	return
}
