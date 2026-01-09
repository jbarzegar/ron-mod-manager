package server

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/mod"
	"github.com/jbarzegar/ron-mod-manager/handler"
	servicev1 "github.com/jbarzegar/ron-mod-manager/proto/service/v1"
	"github.com/jbarzegar/ron-mod-manager/proto/service/v1/servicev1connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func NewGRPCServer(db *ent.Client, h handler.Handler, conf ServerConf) error {
	mux := http.NewServeMux()
	path, handler := servicev1connect.NewServiceHandler(&serviceHandlerServer{})
	mux.Handle(path, handler)

	return http.ListenAndServe(
		conf.Addr,
		h2c.NewHandler(mux, &http2.Server{}),
	)

}

type serviceHandlerServer struct {
	servicev1connect.UnimplementedServiceHandler
	db   *ent.Client
	h    handler.Handler
	conf ServerConf
}

func mapAddArchiveResponse(result *handler.AddModResponse) *servicev1.AddArchiveResponse {
	choices := []*servicev1.Choice{}
	for _, c := range result.Choices {
		choice := &servicev1.Choice{
			Id:       c.UUID.String(),
			Name:     c.Name,
			FullPath: c.FullPath,
		}
		choices = append(choices, choice)
	}

	response := servicev1.AddArchiveResponse{
		Archive: &servicev1.Archive{
			Name:        result.Archive.Name,
			ArchivePath: result.Archive.ArchivePath,
			Md5Sum:      result.Archive.Md5Sum,
			Installed:   result.Archive.Installed,
		},
		Choices: choices,
	}
	return &response
}

func (s *serviceHandlerServer) AddArchive(ctx context.Context,
	req *servicev1.AddArchiveRequest,
) (*servicev1.AddArchiveResponse, error) {
	if req.ArchivePath == "" {
		return nil, errors.New("archive_path must be provided")
	}

	if req.Name == "" {
		return nil, errors.New("name must be provided")
	}

	result, err := s.h.AddArchive(req.ArchivePath, req.Name)
	if err != nil {
		slog.Error("Error adding mod")
		slog.Error(err.Error())
		return nil, err
	}

	return mapAddArchiveResponse(result), nil
}

// mapStateToEnum maps the incoming state string
func mapStateToEnum(s string) servicev1.ModState {
	switch s {
	case "":
		return servicev1.ModState_MOD_STATE_UNSPECIFIED
	case mod.StateActivate.String():
		return servicev1.ModState_MOD_STATE_ACTIVE
	case mod.StateInactive.String():
		return servicev1.ModState_MOD_STATE_INACTIVE
	default:
		// todo: warn
		return servicev1.ModState_MOD_STATE_UNSPECIFIED
	}
}

func (s *serviceHandlerServer) GetAllMods(ctx context.Context,
	req *servicev1.GetAllModsRequest,
) (*servicev1.GetAllModsResponse, error) {
	staged, err := s.h.GetAllMods()
	if err != nil {
		return nil, err
	}

	// map handler -> grpc data
	data := make([]*servicev1.Mod, len(staged.Mods))
	for _, mod := range staged.Mods {
		data = append(data,
			&servicev1.Mod{
				Name:          mod.Name,
				State:         mapStateToEnum(mod.State),
				ActiveVersion: &mod.ActiveVersion,
				Origin:        mod.Origin,
			},
		)
	}
	resp := servicev1.GetAllModsResponse{
		Data: data,
	}

	return &resp, nil
}
