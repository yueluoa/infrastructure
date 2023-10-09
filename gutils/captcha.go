package gutils

import (
	"bytes"
	"github.com/dchest/captcha"
	"net/http"
)

// 当开启多服务器部署时，替换下面的配置，使用redis共享存储验证码
// var store = captcha.NewDefaultRedisStore()
//var store = base64Captcha.DefaultMemStore
//
//type CaptchaResponse struct {
//	ID     string
//	Encode string
//	Err    error
//}
//
//func NewDriver() *base64Captcha.DriverString {
//	driver := new(base64Captcha.DriverString)
//	driver.Height = 44
//	driver.Width = 120
//	driver.NoiseCount = 5
//	driver.ShowLineOptions = base64Captcha.OptionShowSineLine | base64Captcha.OptionShowSlimeLine | base64Captcha.OptionShowHollowLine
//	driver.Length = 6
//	driver.Source = "1234567890qwertyuipkjhgfdsazxcvbnm"
//	driver.Fonts = []string{"wqy-microhei.ttc"}
//	return driver
//}
//
//// 生成图形验证码
//func GenerateCaptchaHandler(w ghttp.ResponseWriter, r *ghttp.Request) {
//	var driver = NewDriver().ConvertFonts()
//	c := base64Captcha.NewCaptcha(driver, store)
//	_, content, answer := c.Driver.GenerateIdQuestionAnswer()
//	id := "captcha:yufei"
//	item, _ := c.Driver.DrawCaptcha(content)
//	c.Store.Set(id, answer)
//	item.WriteTo(w)
//}
//
//func Captcha() *CaptchaResponse {
//	driverString := base64Captcha.DriverString{
//		Height:          50,
//		Width:           120,
//		NoiseCount:      0,
//		ShowLineOptions: 2 | 4,
//		Length:          5,
//		Source:          "1234567890qwertyuioplkjhgfdsazxcvbnmQWERTYUIOPLKJHGFDSAZXCVBNM",
//		//BgColor:         &color.RGBA{R: 3, G: 102, B: 214, A: 125},
//		Fonts: []string{"wqy-microhei.ttc"},
//	}
//	// ConvertFonts 按名称加载字体
//	driver := driverString.ConvertFonts()
//	captcha := base64Captcha.NewCaptcha(driver, store)
//	id, b64s, err := captcha.Generate()
//	if err != nil {
//		return nil
//	}
//	_ = captcha.Store.Set(id, b64s)
//
//	return &CaptchaResponse{
//		ID:     id,
//		Encode: b64s,
//		Err:    err,
//	}
//}
//
//func VerifyCaptcha(id string, VerifyValue string) bool {
//	if store.Verify(id, VerifyValue, true) {
//		return true
//	}
//	return false
//}

const DefaultLen = 4
const DefaultImgWidth = 99
const DefaultImgHeight = 36

func CaptchaHandler(imgWidth, imgHeight int) http.Handler {
	return captcha.Server(imgWidth, imgHeight)
}

func CreateCaptchaImage(length int) string {
	captchaID := captcha.NewLen(length)

	return captchaID
}

func VerifyCaptcha(captchaID, captchaEncode string) bool {
	return captcha.VerifyString(captchaID, captchaEncode)
}

func Reload(captchaID string) bool {
	return captcha.Reload(captchaID)
}

func Verify(captchaID, value string) bool {
	return captcha.VerifyString(captchaID, value)
}

func GetImageByte(captchaID string) []byte {
	var content bytes.Buffer

	err := captcha.WriteImage(&content, captchaID, captcha.StdWidth, captcha.StdHeight)
	if err != nil {
		return nil
	}

	return content.Bytes()
}
