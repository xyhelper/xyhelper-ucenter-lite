package consts

import "github.com/gogf/gf/v2/errors/gcode"

var (
	//CodeNil           = New(200, "OK", "")
	CodeBadRequest    = New(400, "Bad Request", "Bad Request")
	CodeNotAuthorized = New(401, "无权限登录", "无权限登录")
	CodeTokenExpired  = New(401, "token已过期", "token已过期")
	CodeInvalidToken  = New(401, "无效的token", "无效的token")
	CodeInternalError = New(500, "内部错误", "内部错误")
)

type BizCode struct {
	code    int
	message string
	detail  BizCodeDetail
}
type BizCodeDetail struct {
	Code     string
	HttpCode int
}

func (c BizCode) BizDetail() BizCodeDetail {
	return c.detail
}

func (c BizCode) Code() int {
	return c.code
}

func (c BizCode) Message() string {
	return c.message
}

func (c BizCode) Detail() interface{} {
	return c.detail
}

func New(httpCode int, code string, message string) gcode.Code {
	return BizCode{
		code:    httpCode,
		message: message,
		detail: BizCodeDetail{
			Code:     code,
			HttpCode: httpCode,
		},
	}
}
