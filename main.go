package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jbarzegar/ron-mod-manager/appconfig"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/handler"
	"github.com/jbarzegar/ron-mod-manager/handlerio"
	"github.com/jbarzegar/ron-mod-manager/server"
	_ "github.com/mattn/go-sqlite3"
)

func setupDb() (*ent.Client, error) {
	cfg := fmt.Sprintf("file:%v?cache=shared&_fk=1", "test/test.sqlite")

	client, err := ent.Open(
		"sqlite3",
		cfg,
	)
	return client, err
}

func main() {
	db, err := setupDb()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	// pre flight setup
	slog.Info("setting up config")
	if err := appconfig.Setup(); err != nil {
		slog.Error("error setting up config")
		log.Fatal(err)
	}
	slog.Info("appconfig setup")

	appConf, err := appconfig.Read()
	if err != nil {
		slog.Error("error reading app config")
		log.Fatal(err)
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

	slog.Info("Server started")
}
