package middleware

import "github.com/gogf/gf/v2/net/ghttp"

// 添加CORS中间件
func MiddlewareCORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
