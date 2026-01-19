package server

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/handler"
	"github.com/jbarzegar/ron-mod-manager/internal/grpcactionsv1"
	"github.com/jbarzegar/ron-mod-manager/proto/service/v1/servicev1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func setupHandler(appHandler handler.Handler) *grpcactionsv1.ServiceHandlerServer {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}))
	logger.Info("grpc server started")

	return &grpcactionsv1.ServiceHandlerServer{
		Handler: appHandler,
		Logger:  logger,
	}
}

func CreateGRPCServer(_ *ent.Client, appHandler handler.Handler, conf ServerConf) error {
	mux := http.NewServeMux()
	path, serviceHandler := servicev1connect.NewServiceHandler(
		setupHandler(appHandler),
	)
	mux.Handle(path, serviceHandler)

	return http.ListenAndServe(
		conf.Addr,
		h2c.NewHandler(mux, &http2.Server{}),
	)

}
