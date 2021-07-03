package enum

type WorkspaceRefreshType int32

const (
	_ WorkspaceRefreshType = iota
	WorkspaceRefreshTypeStart
	WorkspaceRefreshTypeEnd
	WorkspaceRefreshTypeClose
)
