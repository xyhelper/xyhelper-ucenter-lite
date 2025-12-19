package service

import (
	"backend/config"
	"context"
	"fmt"
	"time"

	"backend/consts"
	"backend/modules/ucenter/model"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

type JwtService struct{}

func NewJwtService() *JwtService {
	return &JwtService{}
}

// 生成accessToken
func (s *JwtService) GenerateToken(ctx context.Context, data *model.TokenData, tokenId string) (string, error) {
	claims := model.Claims{
		UserId:   data.UserId,
		Username: data.Username,
		Nickname: data.Nickname,
		Email:    data.Email,
		Products: data.Products,

		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(g.Cfg().MustGet(ctx, "jwt.expire").Int()) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    g.Cfg().MustGet(ctx, "jwt.issuer").String(),
			Audience:  []string{g.Cfg().MustGet(ctx, "jwt.audience").String()},
		},
	}

	// 使用HMAC签名
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.JwtSecretKey))
	if err != nil {
		return "", err
	}

	// 将tokenId和token存入Redis，key格式为 "user:token:{userId}:{tokenId}"
	redisKey := fmt.Sprintf("application:user:token:%d:%s", data.UserId, tokenId)
	err = config.Redis.SetEX(ctx, redisKey, tokenString, int64(g.Cfg().MustGet(ctx, "jwt.expire").Int()))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 验证Token
func (s *JwtService) VerifyToken(ctx context.Context, tokenString string) (*model.Claims, error) {
	// 1. 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.JwtSecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*model.Claims); ok && token.Valid {
		// 验证Token版本
		if claims.Version != g.Cfg().MustGet(ctx, "jwt.version").String() {
			return nil, gerror.NewCode(consts.CodeNotAuthorized, "token解析错误")
		}
		// 检查Redis中是否存在该token
		redisKey := fmt.Sprintf("application:user:token:%d:%s", claims.UserId, claims.TokenId)
		validToken, err := config.Redis.Get(ctx, redisKey)
		if err != nil {
			return nil, err
		}
		if validToken.String() != tokenString {
			return nil, gerror.NewCode(consts.CodeNotAuthorized, "token已失效")
		}

		return claims, nil
	}
	return nil, gerror.NewCode(consts.CodeNotAuthorized, "token解析错误")
}

// 生成refreshToken
func (s *JwtService) GenerateRefreshToken(ctx context.Context, userId uint, tokenId string) (string, error) {
	claims := model.RefreshClaims{
		UserId:  userId,
		Version: g.Cfg().MustGet(ctx, "jwt.version").String(),
		TokenId: tokenId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(g.Cfg().MustGet(ctx, "jwt.refresh_expire").Int()) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    g.Cfg().MustGet(ctx, "jwt.issuer").String(),
			Audience:  []string{g.Cfg().MustGet(ctx, "jwt.audience").String()},
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(g.Cfg().MustGet(ctx, "jwt.secret").String()))
	if err != nil {
		return "", err
	}
	// 将刷新token存入Redis，key格式为 "user:refresh_token:{userId}:{tokenId}"
	redisKey := fmt.Sprintf("application:user:refresh_token:%d:%s", userId, tokenId)
	err = config.Redis.SetEX(ctx, redisKey, tokenString, int64(g.Cfg().MustGet(ctx, "jwt.refresh_expire").Int()))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// 撤销Token
func (s *JwtService) RevokeToken(ctx context.Context, tokenString string) error {
	claims, err := s.VerifyToken(ctx, tokenString)
	if err != nil {
		return err
	}
	// 从Redis中删除token
	redisKey := fmt.Sprintf("application:user:token:%d:%s", claims.UserId, claims.TokenId)
	_, err = config.Redis.Del(ctx, redisKey)
	if err != nil {
		return err
	}
	// 从Redis中删除对应的refresh token
	refreshTokenKey := fmt.Sprintf("application:user:refresh_token:%d:%s", claims.UserId, claims.TokenId)
	_, err = config.Redis.Del(ctx, refreshTokenKey)
	if err != nil {
		return err
	}
	return nil
}

// 刷新accessToken和refreshToken
func (s *JwtService) RefreshToken(ctx context.Context, refreshToken string) (string, string, error) {
	// 解析刷新Token
	token, err := jwt.ParseWithClaims(refreshToken, &model.RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(g.Cfg().MustGet(ctx, "jwt.secret").String()), nil
	})
	if err != nil {
		return "", "", err
	}
	if claims, ok := token.Claims.(*model.RefreshClaims); ok && token.Valid {
		// 检查Token版本
		if claims.Version != g.Cfg().MustGet(ctx, "jwt.version").String() {
			return "", "", gerror.New("token version mismatch")
		}
		// 检查Redis中是否存在该刷新token
		oldRefreshTokenKey := fmt.Sprintf("application:user:refresh_token:%d:%s", claims.UserId, claims.TokenId)
		validToken, err := config.Redis.Get(ctx, oldRefreshTokenKey)
		if err != nil {
			return "", "", err
		}
		if validToken.String() != refreshToken {
			return "", "", gerror.New("refresh token has been revoked")
		}
		// 生成token需要的数据
		var ucenterUser *model.UcenterUser
		err = cool.DBM(model.NewUcenterUser()).Where("id = ? and status = 1", claims.UserId).Scan(&ucenterUser)
		if err != nil {
			return "", "", err
		}
		tokenData, err := s.GenerateTokenData(ctx, ucenterUser)
		if err != nil {
			return "", "", err
		}
		tokenId := uuid.New().String() // 生成新的tokenId
		// 生成新的accessToken
		accessToken, err := s.GenerateToken(ctx, tokenData, tokenId)
		if err != nil {
			return "", "", err
		}
		// 生成新的refreshToken
		newRefreshToken, err := s.GenerateRefreshToken(ctx, claims.UserId, tokenId)
		if err != nil {
			return "", "", err
		}

		// 删除旧的token
		oldAccessTokenKey := fmt.Sprintf("application:user:token:%d:%s", claims.UserId, claims.TokenId)
		_, err = config.Redis.Del(ctx, oldRefreshTokenKey, oldAccessTokenKey)
		if err != nil {
			// 即使删除失败，也返回新的token，因为新的token已经生成并存入Redis
			g.Log().Error(ctx, "删除旧token失败:", err)
		}

		return accessToken, newRefreshToken, nil
	}
	return "", "", gerror.New("invalid refresh token")
}

// 构建tokenData数据结构
func (s *JwtService) GenerateTokenData(ctx context.Context, ucenterUser *model.UcenterUser) (*model.TokenData, error) {
	// 获取用户权限（服务列表）
	var allServices []g.Map
	if ucenterUser.Permis != "" {
		permisList := gstr.Split(ucenterUser.Permis, ",")
		for _, permis := range permisList {
			permis = gstr.Trim(permis)
			if permis == "" {
				continue
			}
			allServices = append(allServices, g.Map{"code": permis, "name": permis})
		}
	}

	// 构建TokenData
	tokenData := &model.TokenData{
		UserId:   ucenterUser.ID,
		Username: ucenterUser.Name,
		Nickname: ucenterUser.Name,
		Email:    ucenterUser.Email,
		Products: allServices,
	}
	return tokenData, nil
}
