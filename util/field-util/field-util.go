package FieldUtil

const (
	MainTable               = "main"                 // 主表
	CurrentWorkspaceIdField = "current_workspace_id" // 当前工作空间
	CurrentEditId           = "current_edit_soft_id" // 当前编辑软件
	SidebarStatusField      = "sidebar_status"       // 侧边栏状态
	MainBoxShowBtnText      = "main_box_show_text"   // 主盒子_展示按钮提示
	MainBoxShowIcon         = "main_box_show_icon"   // 主盒子_展示图标
	MainBoxShowNum          = "main_box_show_num"    // 主盒子_展示序号
	NoteTextDocument        = "note_text_document"   // 笔记页面_文档_文字提示
	NoteTextMainPoint       = "note_text_main_point" // 笔记页面_大纲_文字提示
	NoteTextDir             = "note_text_dir"        // 笔记页面_目录_文字提示

	LayoutTable     = "layout"          // 布局表( MainTable 的子表)
	WorkspaceTable  = "workspace"       // 工作空间表
	NoteTable       = "note"            // 笔记表
	MarkdownTable   = "markdown"        // Markdown文件表
	ImageCacheTable = "image_cache"     // 图片缓存表
	EditSoftTable   = "edit_soft_cache" // 编辑软件表
)
