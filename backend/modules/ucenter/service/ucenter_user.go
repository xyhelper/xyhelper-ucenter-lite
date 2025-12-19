package service

import (
	"backend/config"
	"backend/modules/ucenter/model"
	"context"

	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
)

type UcenterUserService struct {
	*cool.Service
}

func NewUcenterUserService() *UcenterUserService {
	return &UcenterUserService{
		&cool.Service{
			Model: model.NewUcenterUser(),
			UniqueKey: g.MapStrStr{
				"token": "Token不能重复",
			},
			NotNullKey: g.MapStrStr{
				"name":        "名称不能为空",
				"token":       "Token不能为空",
				"expire_time": "过期时间不能为空",
			},
			PageQueryOp: &cool.QueryOp{
				FieldEQ:      []string{"name", "token", "email", "status"},
				KeyWordField: []string{"name", "token", "email", "status"},
			},
		},
	}
}

// 重写添加方法
func (s *UcenterUserService) ServiceAdd(ctx context.Context, req *cool.AddReq) (data interface{}, err error) {
	var (
		m      = cool.DBM(s.Model)
		r      = g.RequestFromCtx(ctx)
		reqmap = r.GetMap()
	)
	g.Log().Debug(ctx, "添加账号请求参数：", reqmap)

	//要查询用户Token是否存在
	count, err := cool.DBM(model.NewUcenterUser()).Where("token = ?", reqmap["token"]).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, gerror.New("Token已存在，请更换")
	}
	// 设置邮箱
	reqmap["email"] = reqmap["token"].(string) + "@" + config.EmailSuffix

	lastInsertId, err := m.Data(reqmap).InsertAndGetId()
	if err != nil {
		return nil, err
	}
	data = g.Map{"id": lastInsertId}
	return
}

// 重新编辑方法
func (s *UcenterUserService) ServiceUpdate(ctx context.Context, req *cool.UpdateReq) (data interface{}, err error) {
	var (
		m      = cool.DBM(s.Model)
		r      = g.RequestFromCtx(ctx)
		reqmap = r.GetMap()
	)
	g.Log().Debug(ctx, "编辑账号请求参数：", reqmap)
	// 根据id查询账号
	ucenterUser, err := cool.DBM(model.NewUcenterUser()).Where("id = ?", reqmap["id"]).One()
	if err != nil {
		return nil, err
	}
	if ucenterUser == nil {
		return nil, gerror.New("账号不存在")
	}
	//要查询token是否存在，且不是当前账号的
	count, err := cool.DBM(model.NewUcenterUser()).
		Where("token = ? AND id <> ?", reqmap["token"], reqmap["id"]).Count()
	if err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, gerror.New("Token已存在，请更换")
	}
	// 设置邮箱
	reqmap["email"] = reqmap["token"].(string) + "@" + config.EmailSuffix

	// 更新账号
	_, err = m.Data(reqmap).Where("id = ?", reqmap["id"]).Update()
	if err != nil {
		return nil, err
	}
	return reqmap, nil
}
