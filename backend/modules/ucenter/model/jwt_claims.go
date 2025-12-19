package model

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/golang-jwt/jwt/v4"
)

// JWT Claims
type Claims struct {
	// 用户基本信息
	UserId   uint   `json:"userId"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	// 权限（服务）信息
	Products []g.Map `json:"products"`

	Version string `json:"version"`  // JWT版本号
	TokenId string `json:"token_id"` // Token唯一标识

	jwt.RegisteredClaims
}

// 刷新Token Claims
type RefreshClaims struct {
	UserId  uint   `json:"userId"`
	AppCode string `json:"appCode"`
	Version string `json:"version"`
	TokenId string `json:"token_id"` // Token唯一标识
	jwt.RegisteredClaims
}

// TokenData 用于生成Token的数据结构
type TokenData struct {
	// 用户基本信息
	UserId   uint
	Username string
	Nickname string
	Email    string
	// 权限（服务）信息
	Products []g.Map `json:"products"`
}
