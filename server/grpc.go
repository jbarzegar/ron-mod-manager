package server

import (
	"net/http"

	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/handler"
	"github.com/jbarzegar/ron-mod-manager/internal/grpcactionsv1"
	"github.com/jbarzegar/ron-mod-manager/proto/service/v1/servicev1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func CreateGRPCServer(_ *ent.Client, appHandler handler.Handler, conf ServerConf) error {
	mux := http.NewServeMux()
	path, serviceHandler := servicev1connect.NewServiceHandler(
		&grpcactionsv1.ServiceHandlerServer{
			Handler: appHandler,
		},
	)
	mux.Handle(path, serviceHandler)

	return http.ListenAndServe(
		conf.Addr,
		h2c.NewHandler(mux, &http2.Server{}),
	)
}
