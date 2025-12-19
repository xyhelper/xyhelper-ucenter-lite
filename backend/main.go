package main

import (
	_ "backend/internal/packed"

	_ "github.com/cool-team-official/cool-admin-go/contrib/drivers/sqlite"

	_ "backend/modules"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"

	"backend/internal/cmd"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func main() {
	s := g.Server()
	// gres.Dump()
	ctx := gctx.New()
	// 若是其他init方法带参数了，在这里初始化

	//设置静态资源目录
	s.SetServerRoot("resource/public")

	cmd.Main.Run(ctx)
}
