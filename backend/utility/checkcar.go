package utility

import (
	"github.com/cool-team-official/cool-admin-go/cool"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/frame/g"
)

type CarInfo struct {
	Carid    string
	Email    string
	Session  string
	Password string
	Type     string
	//TeamIds  []string
}

func CheckCar(ctx g.Ctx, carid string) (carInfo *CarInfo, err error) {
	// g.Log().Info(ctx, "check carid:", carid)
	sessionVar, err := cool.CacheManager.Get(ctx, "session:"+carid)
	if err != nil {
		return
	}
	sessionJson := gjson.New(sessionVar)

	carInfo = &CarInfo{}
	carInfo.Carid = carid
	carInfo.Email = sessionJson.Get("email").String()
	carInfo.Session = sessionJson.Get("session").String()
	carInfo.Password = sessionJson.Get("password").String()
	carInfo.Type = sessionJson.Get("type").String()
	return
}
