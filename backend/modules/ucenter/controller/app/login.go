package app

import (
	v1 "backend/api/v1"
	"backend/modules/ucenter/service"
	"context"

	"github.com/gogf/gf/v2/net/ghttp"
)

var LoginController = loginController{}

type loginController struct{}

// Authorize 处理 子应用 授权请求
func (a *loginController) Authorize(ctx context.Context, req *v1.AuthorizeReq) (res *v1.AuthorizeRes, err error) {
	// 从 context 中获取 request 对象
	r := ghttp.RequestFromCtx(ctx)
	authURL := "/login.html" //输入Token页面

	r.Session.Set("RedirectUri", req.RedirectUri)
	// 重定向到 login 页面
	r.Response.RedirectTo(authURL)
	return
}

// Login 处理 login页面登录请求
func (a *loginController) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	res, err = service.NewLoginService().Login(ctx, req)
	return
}

// OauthToken处理 code换取token
func (a *loginController) OauthToken(ctx context.Context, req *v1.OauthTokenReq) (res *v1.OauthTokenRes, err error) {
	res, err = service.NewLoginService().OauthToken(ctx, req)
	return res, err
}
