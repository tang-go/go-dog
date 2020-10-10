package api

import (
	"go-dog/cmd/go-dog-ctl/param"
	"go-dog/pkg/log"
	"go-dog/pkg/rand"
	"go-dog/plugins"

	"github.com/mojocn/base64Captcha"
)

//GetCode 验证图片验证码
func (pointer *Service) GetCode(ctx plugins.Context, request param.GetCodeReq) (response param.GetCodeRes, err error) {
	d := base64Captcha.NewDriverString(80, 240, 80, base64Captcha.OptionShowHollowLine, 5, rand.StringRand(6), nil, []string{})
	driver := d.ConvertFonts()
	code := base64Captcha.NewCaptcha(driver, pointer)
	id, b64s, err := code.Generate()
	if err != nil {
		log.Errorln(err.Error())
		return
	}
	response.ID = id
	response.Img = b64s
	return
}
