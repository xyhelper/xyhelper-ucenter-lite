package service

import (
	v1 "backend/api/v1"
	"backend/consts"
	"backend/modules/ucenter/model"
	"context"
	"encoding/json"
	"net/url"
	"time"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/google/uuid"
)

type LoginService struct {
}

func NewLoginService() *LoginService {
	return &LoginService{}
}

func (s *LoginService) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	r := ghttp.RequestFromCtx(ctx)

	userToken := req.UserToken
	if userToken == "" {
		r.Response.RedirectTo("/login.html?error=" + url.QueryEscape("userToken不能为空"))
		return
	}
	//g.Log().Debug(ctx, " userToken:", userToken)

	// 根据userToken查询ucenter_user用户信息
	var ucenterUser *model.UcenterUser
	err = cool.DBM(model.NewUcenterUser()).Where("token = ? and status = 1 and expire_time > ?", userToken, time.Now()).Scan(&ucenterUser)
	if err != nil {
		r.Response.RedirectTo("/login.html?error=" + url.QueryEscape("查询userToken异常："+err.Error()))
		return
	}
	if ucenterUser == nil {
		r.Response.RedirectTo("/login.html?error=" + url.QueryEscape("无效的userToken"))
		return
	}
	// 生成token数据
	tokenData, err := NewJwtService().GenerateTokenData(ctx, ucenterUser)
	if err != nil {
		r.Response.RedirectTo("/login.html?error=" + url.QueryEscape("生成token数据异常："+err.Error()))
		return
	}

	tokenId := uuid.New().String() // 生成新的tokenId
	// 生成accessToken
	accessToken, err := NewJwtService().GenerateToken(ctx, tokenData, tokenId)
	if err != nil {
		return nil, err
	}
	// 生成refreshToken
	refreshToken, err := NewJwtService().GenerateRefreshToken(ctx, ucenterUser.ID, tokenId)
	if err != nil {
		return nil, err
	}
	expiresIn := g.Cfg().MustGet(ctx, "jwt.expire").Int()
	refreshExpiresIn := g.Cfg().MustGet(ctx, "jwt.refresh_expire").Int()

	// 生成随机code，将code与token信息放入redis中
	subCode := grand.S(32) // 生成32位随机字符串作为code
	tokenJson, _ := json.Marshal(map[string]interface{}{
		"access_token":       accessToken,
		"refresh_token":      refreshToken,
		"expires_in":         expiresIn,
		"refresh_expires_in": refreshExpiresIn,
		"token_type":         "Bearer",
	})

	clearsDuration := 1 * time.Hour
	err = cool.CacheManager.Set(ctx, "oauth:code:"+subCode, string(tokenJson), clearsDuration)
	if err != nil {
		return nil, gerror.NewCode(consts.CodeInternalError, "存储code失败："+err.Error())
	}

	// 从session中获取回调地址，拼接上code
	// r := ghttp.RequestFromCtx(ctx)
	redirectUri := r.Session.MustGet("RedirectUri").String()
	redirectUri = redirectUri + "?code=" + subCode
	g.Log().Debug(ctx, "重定向地址：", redirectUri)

	//重定向改地址
	r.Response.RedirectTo(redirectUri)
	return
}

func (s *LoginService) OauthToken(ctx context.Context, req *v1.OauthTokenReq) (res *v1.OauthTokenRes, err error) {
	//g.Log().Debug(ctx, "子应用code换取token请求参数：", req)
	r := g.RequestFromCtx(ctx)

	if req.GrantType == "" || req.GrantType == "authorization_code" { // 默认使用authorization_code
		code := req.Code
		if code == "" {
			r.Response.WriteStatus(500)
			r.Response.WriteJson(g.Map{
				"message": "code不能为空",
			})
			return nil, nil
		}
		//g.Log().Debug(ctx, "换取accessToken的code:", code)

		// 从redis中获取token信息
		tokenJsonInfo, err := cool.CacheManager.Get(ctx, "oauth:code:"+code)
		if err != nil {
			r.Response.WriteStatus(400)
			r.Response.WriteJson(g.Map{
				"message": "获取code失败：" + err.Error(),
			})
			return nil, nil
		}
		tokenJsonStr := tokenJsonInfo.String()
		if tokenJsonStr == "" {
			r.Response.WriteStatus(400)
			r.Response.WriteJson(g.Map{
				"message": "code无效或已过期",
			})
			return nil, nil
		}
		tokenJson := gjson.New(tokenJsonStr)
		//删除 该code的redis缓存
		_, err = cool.CacheManager.Remove(ctx, "oauth:code:"+code)
		if err != nil {
			//return nil, gerror.NewCode(consts.CodeInternalError, "删除code缓存失败："+err.Error())
			r.Response.WriteStatus(500)
			r.Response.WriteJson(g.Map{
				"message": "删除code缓存失败：" + err.Error(),
			})
			return nil, nil
		}
		g.RequestFromCtx(ctx).Response.WriteJson(g.Map{
			"access_token":       tokenJson.Get("access_token").String(),
			"refresh_token":      tokenJson.Get("refresh_token").String(),
			"expires_in":         int(tokenJson.Get("expires_in").Int()),
			"refresh_expires_in": int(tokenJson.Get("refresh_expires_in").Int()),
			"token_type":         tokenJson.Get("token_type").String(),
		})
	} else if req.GrantType == "refresh_token" {
		refreshToken := req.RefreshToken
		if refreshToken == "" {
			return nil, gerror.NewCode(consts.CodeNotAuthorized, "未传递需要的refreshToken")
		}
		//刷新token
		accessToken, refreshToken, err := NewJwtService().RefreshToken(ctx, refreshToken)
		if err != nil {
			return nil, err
		}
		expiresIn := g.Cfg().MustGet(ctx, "jwt.expire").Int()
		refreshExpiresIn := g.Cfg().MustGet(ctx, "jwt.refresh_expire").Int()

		g.RequestFromCtx(ctx).Response.WriteJson(g.Map{
			"access_token":       accessToken,
			"refresh_token":      refreshToken,
			"expires_in":         expiresIn,
			"refresh_expires_in": refreshExpiresIn,
			"token_type":         "Bearer",
		})
	}
	return nil, nil
}
