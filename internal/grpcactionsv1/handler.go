package grpcactionsv1

import (
	"context"
	"errors"
	"log/slog"

	"github.com/jbarzegar/ron-mod-manager/handler"
	"github.com/jbarzegar/ron-mod-manager/internal/actions"
	servicev1 "github.com/jbarzegar/ron-mod-manager/proto/service/v1"
	"github.com/jbarzegar/ron-mod-manager/proto/service/v1/servicev1connect"
)

type ServiceHandlerServer struct {
	servicev1connect.UnimplementedServiceHandler
	Handler handler.Handler
}

func (s *ServiceHandlerServer) AddArchive(ctx context.Context,
	req *servicev1.AddArchiveRequest,
) (*servicev1.AddArchiveResponse, error) {
	if req.ArchivePath == "" {
		return nil, errors.New("archive_path must be provided")
	}

	if req.Name == "" {
		return nil, errors.New("name must be provided")
	}

	result, err := s.Handler.AddArchive(req.ArchivePath, req.Name)
	if err != nil {
		slog.Error("Error adding mod")
		slog.Error(err.Error())
		return nil, err
	}

	return mapAddArchiveResponse(result), nil
}

func (s *ServiceHandlerServer) GetAllMods(ctx context.Context,
	req *servicev1.GetAllModsRequest,
) (*servicev1.GetAllModsResponse, error) {
	result, err := s.Handler.GetAllMods()
	if err != nil {
		return nil, err
	}

	resp := mapGetAllModsResponse(result)
	return &resp, nil
}

func (s *ServiceHandlerServer) GetStagedMods(ctx context.Context,
	req *servicev1.GetStagedModsRequest,
) (*servicev1.GetStagedModsResponse, error) {
	return nil, errors.New("not implemented")
}

func (s *ServiceHandlerServer) InstallMod(ctx context.Context,
	req *servicev1.InstallModRequest,
) (*servicev1.InstallModResponse, error) {
	modID, modVersionUUID, choices, err := mapInstallModRequest(req)
	if err != nil {
		return nil, err
	}

	if err := s.Handler.InstallMod(modID, modVersionUUID, choices); err != nil {
		return nil, err
	}

	return &servicev1.InstallModResponse{}, nil
}

func (s *ServiceHandlerServer) UninstallMod(ctx context.Context,
	req *servicev1.UninstallModRequest,
) (*servicev1.UninstallModResponse, error) {
	modIds := make([]int, len(req.GetModIds()))

	for _, mid := range req.GetModIds() {
		modIds = append(modIds, int(mid))
	}

	if err := s.Handler.UninstallMod(modIds); err != nil {
		return nil, err
	}

	return &servicev1.UninstallModResponse{}, nil
}

func (s *ServiceHandlerServer) DeleteMod(ctx context.Context,
	req *servicev1.DeleteModRequest,
) (*servicev1.DeleteModResponse, error) {
	modsToDelete := make([]actions.DeleteModEntry, len(req.GetModsToDelete()))

	for _, td := range req.GetModsToDelete() {
		modsToDelete = append(modsToDelete, actions.DeleteModEntry{
			ID:             int(td.GetId()),
			DeleteVersions: mapDeleteModType(td.GetDeleteStates()),
			DeleteArchive:  td.GetDeleteArchive(),
		})
	}

	if err := s.Handler.DeleteMod(actions.DeleteModRequest{Mods: modsToDelete}); err != nil {
		return nil, err
	}

	return &servicev1.DeleteModResponse{}, nil
}
