package api

import (
	"go-dog/cmd/define"
	"go-dog/cmd/go-dog-ctl/table"
	customerror "go-dog/error"
	"go-dog/plugins"
)

//Auth 插件
func (pointer *API) Auth(ctx plugins.Context, token string) error {
	admin := new(table.Admin)
	if e := pointer.cache.GetCache().Get(token, admin); e != nil {
		return customerror.EnCodeError(define.AdminTokenErr, "token失效或者不正确")
	}
	ctx.SetShare("Admin", admin)
	return nil
}
