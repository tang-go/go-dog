package uuid

import (
	"fmt"

	uuid "github.com/satori/go.uuid"
)

//GetToken 获取token
func GetToken() string {
	token, err := uuid.NewV4()
	if err != nil {
		panic(err.Error())
	}
	t := fmt.Sprintf("%s", token)
	return t
}
