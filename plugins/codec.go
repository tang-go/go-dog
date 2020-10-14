package plugins

//Codec 参数编码器
type Codec interface {

	//EnCode 编码
	EnCode(code string, v interface{}) ([]byte, error)

	//DeCode 编码
	DeCode(code string, buff []byte, v interface{}) error
}
