package plugins

//Auth 验证
type Auth interface {
	//Auth 验证
	Auth(ctx Context, token string) error
}
