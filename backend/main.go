/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	connect "connectrpc.com/connect"
	"github.com/jbarzegar/ron-mod-manager/db"
	"github.com/jbarzegar/ron-mod-manager/rpc"
	"github.com/jbarzegar/ron-mod-manager/rpc/rpcconnect"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

const address = "localhost:8080"

func main() {
	client, err := db.InitClient()
	if err != nil {
		log.Fatalln(err)
	}

	defer client.Close()

	mux := http.NewServeMux()
	path, handler := rpcconnect.NewRPCHandler(&rpcServiceServer{})
	mux.Handle(path, handler)

	fmt.Println("Listening on", address)
	http.ListenAndServe(address, h2c.NewHandler(mux, &http2.Server{}))

	// TODO: Readd the following =>
	//
	// if err := client.Schema.Create(context.Background()); err != nil {
	// 	log.Fatalf("failed creating schema resources: %v", err)
	// }

	// stateManagement.PreflightChecks()
	// cmd.Execute()
}

type rpcServiceServer struct {
	rpcconnect.UnimplementedRPCHandler
}

func (s *rpcServiceServer) InstallArchive(ctx context.Context, req *connect.Request[rpc.InstallArchiveRequest]) (*connect.Response[rpc.InstallArchiveReply], error) {
	fmt.Println(req.Msg.GetModPath())

	return connect.NewResponse(&rpc.InstallArchiveReply{Status: false}), nil
}
