package api

import (
	"fmt"
	"go-dog/cmd/go-dog-ctl/param"
	"go-dog/plugins"
)

//GetAPIList 获取api列表
func (pointer *API) GetAPIList(ctx plugins.Context, req param.GetAPIListReq) (res param.GetAPIListRes, err error) {
	fmt.Println(ctx.GetShareByKey("Admin"))
	list := make(map[string]*param.Service)
	pointer.lock.RLock()
	for key, api := range pointer.apis {
		if service, ok := list[api.name]; ok {
			service.APIS = append(service.APIS, &param.API{
				Name:     api.method.Name,
				Level:    api.method.Level,
				Request:  api.method.Request,
				Response: api.method.Response,
				Explain:  api.method.Explain,
				IsAuth:   api.method.IsAuth,
				Version:  api.method.Version,
				URL:      key,
				Kind:     api.method.Kind,
			})
		} else {
			s := &param.Service{
				Name:    api.name,
				Explain: api.explain,
				APIS: []*param.API{
					&param.API{
						Name:     api.method.Name,
						Level:    api.method.Level,
						Request:  api.method.Request,
						Response: api.method.Response,
						Explain:  api.method.Explain,
						IsAuth:   api.method.IsAuth,
						Version:  api.method.Version,
						URL:      key,
						Kind:     api.method.Kind,
					},
				},
			}
			list[api.name] = s
		}
	}
	pointer.lock.RUnlock()
	for _, s := range list {
		res.List = append(res.List, s)
	}
	return
}
