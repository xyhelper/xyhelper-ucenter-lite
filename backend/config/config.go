package config

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	EmailSuffix  = "outlook.com" // 默认邮箱后缀
	JwtSecretKey = ""            // 子应用下发token秘钥

	Redis = g.Redis("cool")
)

func init() {
	ctx := gctx.GetInitCtx()
	// 子应用jwt秘钥
	jwtSecretKey := g.Cfg().MustGetWithEnv(ctx, "JWT_SECRET_KEY").String()
	if jwtSecretKey != "" {
		JwtSecretKey = jwtSecretKey
	}
	// 邮箱后缀
	emailSuffix := g.Cfg().MustGetWithEnv(ctx, "EMAIL_SUFFIX").String()
	if emailSuffix != "" {
		EmailSuffix = emailSuffix
	}

}
