package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"kxtGo/controller"
	"kxtGo/tool"
)

func  main()  {
	app := gin.Default()

	registerRouter(app)

	//加载app配置
	cfg := new(tool.Config)
	appCfg,err := cfg.GetSection("app")
	if err != nil {
		panic("app配置信息错误")
	}

	// 启动服务
	if err := app.Run(appCfg["host"]+":"+appCfg["port"]); err != nil{
		fmt.Println("err>>> ", err.Error())
	}
}

//注册路由
func registerRouter(router *gin.Engine)  {
	Group := router.Group("/trademark")
	{
		new(controller.ExcelImportTrademarkController).Router(Group)
		new(controller.ImagesTypeTrademarkController).Router(Group)
	}
}