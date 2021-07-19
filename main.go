package main

import (
	"github.com/gin-gonic/gin"
	LayoutController "yiu/yiu-reader/controller/layout-controller"
	MainController "yiu/yiu-reader/controller/main-controller"
	NoteController "yiu/yiu-reader/controller/note-controller"
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
	router.PUT("/current/workspace/:id", MainController.SetCurrentWorkspace)

	workspaceGroup := router.Group("/workspace")
	{
		workspaceGroup.POST("", WorkspaceController.Add)
		workspaceGroup.GET("", WorkspaceController.Search)
		workspaceGroup.GET("/:id", WorkspaceController.View)
		workspaceGroup.PUT("", WorkspaceController.Update)
		workspaceGroup.DELETE("/:id", WorkspaceController.Delete)
		workspaceGroup.PUT("/up/:id", WorkspaceController.Up)
		workspaceGroup.PUT("/down/:id", WorkspaceController.Down)
		// workspaceGroup.GET("/content", WorkspaceController.Content)
	}

	layoutGroup := router.Group("/layout")
	{
		layoutGroup.GET("", LayoutController.Search)
		layoutGroup.POST("", LayoutController.Add)
		layoutGroup.PUT("/resize", LayoutController.ResizePosition)
		layoutGroup.DELETE("/:id", LayoutController.Delete)
		layoutGroup.PUT("", LayoutController.Update)
		layoutGroup.GET("/:id", LayoutController.View)
	}

	noteGroup := router.Group("/note")
	{
		noteGroup.GET("/refresh", NoteController.Refresh)
		noteGroup.POST("/tree", NoteController.SearchTree)
		noteGroup.GET("", NoteController.Search)
		noteGroup.DELETE("/:id", NoteController.Delete)
		noteGroup.DELETE("/file/:id", NoteController.DeleteFile)
		noteGroup.GET("/position/:id", NoteController.Position)
		noteGroup.GET("/change/show/:id", NoteController.ChangeShow)
		noteGroup.PUT("/up/:id", NoteController.Up)
		noteGroup.PUT("/down/:id", NoteController.Down)
	}

	_ = router.Run(":8080")
}
