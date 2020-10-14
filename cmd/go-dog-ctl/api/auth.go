package api

import "go-dog/plugins"

//Auth 验证插件
func (pointer *Service) Auth(ctx plugins.Context, token string) error {
	_, err := pointer._Auth(token)
	if err != nil {
		return err
	}
	return nil
}
