package cmd

import (
	"context"

	"backend/modules/ucenter/controller/app"
	"backend/modules/ucenter/middleware"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
)

var (
	Main = gcmd.Command{
		Name:  "main",
		Usage: "main",
		Brief: "start http server",
		Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
			if cool.IsRedisMode {
				go cool.ListenFunc(ctx)
			}
			s := g.Server()

			s.Group("/", func(group *ghttp.RouterGroup) {
				group.Middleware(middleware.MiddlewareCORS) // 添加CORS中间件
				group.Bind(
					app.LoginController.Authorize,  //子应用登录（接口地址：/authorize）
					app.LoginController.OauthToken, //子应用登录Code换取token(接口地址：/oauth/token)

					app.LoginController.Login, //login页面登录请求（接口地址：/login）
				)
			})

			s.Run()
			return nil
		},
	}
)
