package plugins

//Codec 参数编码器
type Codec interface {

	//EnCode 编码
	EnCode(v interface{}) ([]byte, error)

	//DeCode 编码
	DeCode(buff []byte, v interface{}) error
}
