/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"log"

	"github.com/jbarzegar/ron-mod-manager/cmd"
	"github.com/jbarzegar/ron-mod-manager/db"
	stateManagement "github.com/jbarzegar/ron-mod-manager/state-management"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	client, err := db.InitClient()

	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	if err := client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	stateManagement.PreflightChecks()
	cmd.Execute()
}
