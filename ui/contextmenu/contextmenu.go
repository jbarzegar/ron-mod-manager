package contextmenu

import "context"

type ContextMenu struct {
	DebugLog *DebugLog
}

func NewContextMenu(ctx context.Context) *ContextMenu {
	return &ContextMenu{DebugLog: NewDebugLog(ctx)}
}
