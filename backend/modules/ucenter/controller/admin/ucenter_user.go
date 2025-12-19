package admin

import (
	"backend/modules/ucenter/service"

	"github.com/cool-team-official/cool-admin-go/cool"
)

type UcenterUserController struct {
	*cool.Controller
}

func init() {
	var ucenter_user_controller = &UcenterUserController{
		&cool.Controller{
			Prefix:  "/admin/ucenter/user",
			Api:     []string{"Add", "Delete", "Update", "Info", "List", "Page"},
			Service: service.NewUcenterUserService(),
		},
	}
	// 注册路由
	cool.RegisterController(ucenter_user_controller)
}
