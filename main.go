package main

import (
	"github.com/gin-gonic/gin"
	MainController "yiu/yiu-reader/controller/main-controller"
	OpUtil "yiu/yiu-reader/util/op-util"
)

func main() {
	// 创建DB Bean
	OpUtil.CreateDB(".yiu/yiu-reader.db")
	// 关闭DBBean
	defer OpUtil.CloseDB()
	// 创建Logger Bean
	OpUtil.CreateLogger()
	// 默认路由引擎
	router := gin.Default()
	// 加载静态文件
	router.Static("/assets", "./dist/assets")
	// index.html
	router.LoadHTMLFiles("dist/index.html")
	router.GET("/", MainController.IndexHTML)

	apiGroup := router.Group("/api")
	{
		mainGroup := apiGroup.Group("/main")
		{
			mainGroup.GET("/current/workspace", MainController.GetCurrentWorkspace)
		}
	}

	_ = router.Run(":8081")
}
