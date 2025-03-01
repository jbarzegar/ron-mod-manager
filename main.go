package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/jbarzegar/ron-mod-manager/appconfig"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/server"
	_ "github.com/mattn/go-sqlite3"
)

func setupDb() (*ent.Client, error) {
	cfg := fmt.Sprintf("file:%v?cache=shared&_fk=1", "test.sqlite")

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
	appconfig.Setup()
	//
	//
	//
	//
	// start server
	if _, err := server.CreateServer(db); err != nil {
		log.Fatal(err)
	} else {
		slog.Info("Server started")
	}
}
