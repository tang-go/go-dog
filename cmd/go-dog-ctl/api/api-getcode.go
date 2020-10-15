package api

import (
	"go-dog/cmd/define"
	"go-dog/cmd/go-dog-ctl/param"
	customerror "go-dog/error"
	"go-dog/lib/rand"
	"go-dog/log"
	"go-dog/plugins"

	"github.com/mojocn/base64Captcha"
)

//GetCode 验证图片验证码
func (pointer *API) GetCode(ctx plugins.Context, request param.GetCodeReq) (response param.GetCodeRes, err error) {
	//查看是否是测试接口
	if ctx.GetIsTest() {
		response.ID = request.Code
		response.Img = "123456"
		pointer.Set(response.ID, response.Img)
		return
	}
	number := rand.StringRand(6)
	d := base64Captcha.NewDriverString(80, 240, 80, base64Captcha.OptionShowHollowLine, 6, number, nil, []string{})
	driver := d.ConvertFonts()
	code := base64Captcha.NewCaptcha(driver, pointer)
	id, b64s, err := code.Generate()
	if err != nil {
		log.Errorln(err.Error())
		err = customerror.EnCodeError(define.GetCodeErr, err.Error())
		return
	}
	response.ID = id
	response.Img = b64s
	return
}
