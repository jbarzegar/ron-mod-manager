package ui

import (
	"context"

	"github.com/jbarzegar/ron-mod-manager/ui/contextmenu"
)

// App struct
type App struct {
	ctx         context.Context
	ContextMenu *contextmenu.ContextMenu
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.ContextMenu = contextmenu.NewContextMenu(ctx)
}

func (a *App) SetupLogs() []contextmenu.LogEntry {
	return a.ContextMenu.DebugLog.GetInitialLogs()
}
