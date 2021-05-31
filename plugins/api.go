package plugins

//API API接口
type API interface {

	//APIGroup APi组
	APIGroup(group string) API

	//APIAuth APi需要验证
	APIAuth() API

	//APINoAuth APi需要不验证
	APINoAuth() API

	//APIVersion APi版本
	APIVersion(version string) API

	//APILevel APi等级
	APILevel(level int8) API

	//GET APi GET路由
	GET(name string, path string, explain string, fn interface{})

	//POST POST路由
	POST(name string, path string, explain string, fn interface{})

	//PUT PUT路由
	PUT(name string, path string, explain string, fn interface{})

	//DELETE DELETE路由
	DELETE(name string, path string, explain string, fn interface{})
}
