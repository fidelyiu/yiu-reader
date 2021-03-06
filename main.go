package main

import (
	yiuLog "github.com/fidelyiu/yiu-go-tool/log"
	"github.com/gin-gonic/gin"
	DbController "yiu/yiu-reader/controller/db-controller"
	EditSoftController "yiu/yiu-reader/controller/edit-soft-controller"
	ImageController "yiu/yiu-reader/controller/image-controller"
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
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	// 加载静态文件
	router.Static("/assets", "./dist/assets")
	err := OpUtil.CreateImageDir()
	if err != nil {
		yiuLog.ErrorLn(err)
		return
	}
	router.Static("/image", "./.yiu/image")
	router.LoadHTMLFiles("dist/index.html")

	// index.html
	router.GET("/", MainController.IndexHTML)

	mainGroup := router.Group("/main")
	{
		mainGroup.GET("/current/workspace", MainController.GetCurrentWorkspace)
		mainGroup.PUT("/current/workspace/:id", MainController.SetCurrentWorkspace)
		mainGroup.GET("/main/box/txt", MainController.GetMainBoxShowText)
		mainGroup.PUT("/main/box/txt", MainController.SetMainBoxShowText)
		mainGroup.GET("/main/box/icon", MainController.GetMainBoxShowIcon)
		mainGroup.PUT("/main/box/icon", MainController.SetMainBoxShowIcon)
		mainGroup.GET("/main/box/num", MainController.GetMainBoxShowNum)
		mainGroup.PUT("/main/box/num", MainController.SetMainBoxShowNum)
		mainGroup.GET("/sidebar/status", MainController.GetSidebarStatus)
		mainGroup.PUT("/sidebar/status", MainController.SetSidebarStatus)
		mainGroup.GET("/edit/soft", MainController.GetEditSoft)
		mainGroup.PUT("/edit/soft/:id", MainController.SetEditSoft)
		mainGroup.GET("/os/pathSeparator", MainController.GetOsPathSeparator)
		mainGroup.GET("/note/txt/document", MainController.GetNoteTextDocument)
		mainGroup.PUT("/note/txt/document", MainController.SetNoteTextDocument)
		mainGroup.GET("/note/txt/main/point", MainController.GetNoteTextMainPoint)
		mainGroup.PUT("/note/txt/main/point", MainController.SetNoteTextMainPoint)
		mainGroup.GET("/note/txt/dir", MainController.GetNoteTextDir)
		mainGroup.PUT("/note/txt/dir", MainController.SetNoteTextDir)
	}

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
		noteGroup.POST("", NoteController.Add)
		noteGroup.PUT("", NoteController.Update)
		noteGroup.GET("/:id", NoteController.View)
		noteGroup.GET("/refresh", NoteController.Refresh)
		noteGroup.POST("/tree", NoteController.SearchTree)
		noteGroup.GET("", NoteController.Search)
		noteGroup.DELETE("/:id", NoteController.Delete)
		noteGroup.DELETE("/file/:id", NoteController.DeleteFile)
		noteGroup.DELETE("/bad/:id", NoteController.DeleteBad)
		noteGroup.GET("/position/:id", NoteController.Position)
		noteGroup.GET("/change/show/:id", NoteController.ChangeShow)
		noteGroup.PUT("/up/:id", NoteController.Up)
		noteGroup.PUT("/down/:id", NoteController.Down)
		noteGroup.GET("/edit/md/:id", NoteController.EditMarkdown)
		noteGroup.GET("/reade/:id", NoteController.Reade)
		noteGroup.GET("/dir/tree/:id", NoteController.DirTree)
		noteGroup.GET("/num/document/:id", NoteController.GetNumDocument)
		noteGroup.PUT("/num/document", NoteController.SetNumDocument)
		noteGroup.GET("/num/main/point/:id", NoteController.GetNumMainPoint)
		noteGroup.PUT("/num/main/point", NoteController.SetNumMainPoint)
		noteGroup.GET("/num/dir/:id", NoteController.GetNumDir)
		noteGroup.PUT("/num/dir", NoteController.SetNumDir)

	}

	editSoftGroup := router.Group("/edit/soft")
	{
		editSoftGroup.POST("", EditSoftController.Add)
		editSoftGroup.GET("", EditSoftController.Search)
		editSoftGroup.GET("/:id", EditSoftController.View)
		editSoftGroup.PUT("", EditSoftController.Update)
		editSoftGroup.DELETE("/:id", EditSoftController.Delete)
		editSoftGroup.PUT("/up/:id", EditSoftController.Up)
		editSoftGroup.PUT("/down/:id", EditSoftController.Down)
	}

	dbGroup := router.Group("/db")
	{
		dbGroup.POST("", DbController.Search)
	}

	imgGroup := router.Group("/img")
	{
		imgGroup.POST("/upload", ImageController.Upload)
		imgGroup.GET("/load", ImageController.Load)
	}

	_ = router.Run(":8080")
}
