package db

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jbarzegar/ron-mod-manager/config"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/archive"
)

var _Client *ent.Client = nil

func Client() *ent.Client {
	if _Client == nil {
		log.Fatalf("Client is not initialized")
		return nil
	}

	return _Client
}

func InitClient() (*ent.Client, error) {
	if _Client != nil {
		return nil, errors.New("client-already-connected")
	}

	cfgStr := fmt.Sprintf("file:%v?cache=shared&_fk=1", config.DBPath)
	client, err := ent.Open("sqlite3", cfgStr)

	if err != nil {
		log.Fatalf("Failed to connect to sqlite", err)
	}

	_Client = client

	return client, nil

}

func GetArchiveByName(name string) (*ent.Archive, error) {
	return Client().Archive.Query().Where(archive.Name(name)).Only(context.Background())
}
