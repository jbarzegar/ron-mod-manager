package ui

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/jbarzegar/ron-mod-manager/appconfig"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/handler"
	"github.com/jbarzegar/ron-mod-manager/handlerio"
	"github.com/jbarzegar/ron-mod-manager/server"
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

func setupDb() (*ent.Client, error) {
	cfg := fmt.Sprintf("file:%v?cache=shared&_fk=1", "test/test.sqlite")

	client, err := ent.Open(
		"sqlite3",
		cfg,
	)
	return client, err
}

func startServer(logger *slog.Logger) {
	logger.Info("Starting up Setup process")

	logger.Info("Setting up DB")
	db, err := setupDb()
	if err != nil {
		logger.Error("failed to setup db", "err", err)
	}

	// pre flight setup
	logger.Info("setting up config")
	if err := appconfig.Setup(); err != nil {
		logger.Error("error setting up config")
		// log.Fatal(err)
	}
	logger.Info("appconfig setup")

	appConf, err := appconfig.Read()
	if err != nil {
		logger.Error("error reading app config")
		// log.Fatal(err)
	}

	// setup handlers for transport layer
	IOHandler := &handlerio.FileSystemHandler{
		Config: appConf,
	}
	appHandler := handler.Handler{Db: db, Config: appConf, Io: IOHandler}

	// start server
	if err := server.CreateGRPCServer(
		db,
		appHandler,
		server.ServerConf{Addr: ":5000"},
	); err != nil {
		log.Fatal(err)
		return
	}

	logger.Info("Server started")

}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.ContextMenu = contextmenu.NewContextMenu(ctx)

	logger := a.ContextMenu.DebugLog.NewLogger()
	startServer(logger)
}

// called in frontend
func (a *App) SetupLogs() []contextmenu.LogEntry {
	return a.ContextMenu.DebugLog.GetInitialLogs()
}

// called in frontend
func (a *App) AddNewLog() {
	a.ContextMenu.DebugLog.NewLogs([]string{
		"new log from event. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
	})

}
