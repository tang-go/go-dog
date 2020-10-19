package api

import (
	"go-dog/cmd/go-dog-ctl/param"
	"go-dog/plugins"
)

//GetServiceList 获取服务列表
func (pointer *API) GetServiceList(ctx plugins.Context, req param.GetServiceReq) (res param.GetServiceRes, err error) {
	services := pointer.service.GetClient().GetAllService()
	for _, service := range services {
		s := &param.ServiceInfo{
			Key:       service.Key,
			Name:      service.Name,
			Address:   service.Address,
			Port:      service.Port,
			Explain:   service.Explain,
			Longitude: service.Longitude,
			Latitude:  service.Latitude,
			Time:      service.Time,
		}
		for _, method := range service.Methods {
			s.Methods = append(s.Methods, &param.Method{
				Name:     method.Name,
				Level:    method.Level,
				Request:  method.Request,
				Response: method.Response,
				Explain:  method.Explain,
				IsAuth:   method.IsAuth,
			})
		}
		res.List = append(res.List, s)
	}
	return
}
