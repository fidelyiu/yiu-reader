package main

import (
	"github.com/gin-gonic/gin"
	MainController "yiu/yiu-reader/controller/main-controller"
	WorkspaceController "yiu/yiu-reader/controller/workspace-controller"
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

	router.GET("/current/workspace", MainController.GetCurrentWorkspace)
	router.PUT("/current/workspace", MainController.SetCurrentWorkspace)

	workspaceGroup := router.Group("/workspace")
	{
		workspaceGroup.POST("", WorkspaceController.Add)
		workspaceGroup.GET("", WorkspaceController.Search)
		workspaceGroup.GET("/:id", WorkspaceController.View)
		workspaceGroup.PUT("", WorkspaceController.Update)
		workspaceGroup.DELETE("/:id", WorkspaceController.Delete)
		// workspaceGroup.GET("/content", WorkspaceController.Content)
	}

	_ = router.Run(":8081")
}
