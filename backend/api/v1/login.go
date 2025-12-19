package v1

import "github.com/gogf/gf/v2/frame/g"

// app 子应用登录请求
type AuthorizeReq struct {
	g.Meta      `path:"/authorize" method:"get" tags:"登录管理" summary:"登录请求"`
	RedirectUri string `json:"redirectUri" v:"required#redirectUri不能为空" dc:"子应用回调地址"`
}

// app 子应用登录响应
type AuthorizeRes struct {
	g.Meta `mime:"application/json"`
}

// app login页面登录请求
type LoginReq struct {
	g.Meta    `path:"/login" method:"get" tags:"登录管理" summary:"login页面登录请求"`
	UserToken string `json:"userToken" v:"required#userToken不能为空" dc:"登录userToken"`
}

// app login页面登录响应
type LoginRes struct {
	g.Meta `mime:"application/json"`
}

// app 子应用code换取token
type OauthTokenReq struct {
	g.Meta       `path:"/oauth/token" method:"post" tags:"登录管理" summary:"授权码换取token"`
	GrantType    string `json:"grant_type" v:"required#grant_type不能为空" dc:"授权类型,authorization_code或者refresh_token，默认authorization_code"`
	Code         string `json:"code" dc:"OAuth授权码"`
	RefreshToken string `json:"refresh_token" dc:"刷新令牌"`
}

// app 子应用code换取token响应
type OauthTokenRes struct {
	g.Meta           `mime:"application/json"`
	AccessToken      string `json:"accessToken" dc:"accessToken"`
	IDToken          string `json:"idToken" dc:"ID Token"`
	RefreshToken     string `json:"refreshToken" dc:"refreshToken"`
	ExpiresIn        int    `json:"expiresIn" dc:"过期时间"`
	RefreshExpiresIn int    `json:"refreshExpiresIn" dc:"刷新过期时间"`
	TokenType        string `json:"tokenType" dc:"token类型"`
}
