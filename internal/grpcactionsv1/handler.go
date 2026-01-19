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
	Logger  *slog.Logger
}

func (s *ServiceHandlerServer) AddArchive(ctx context.Context,
	req *servicev1.AddArchiveRequest,
) (*servicev1.AddArchiveResponse, error) {
	s.Logger.Info("AddArchive")
	if req.ArchivePath == "" {
		s.Logger.Error("no archive provided")
		return nil, errors.New("archive_path must be provided")
	}

	if req.Name == "" {
		s.Logger.Error("no name provided")
		return nil, errors.New("name must be provided")
	}

	result, err := s.Handler.AddArchive(req.ArchivePath, req.Name)
	if err != nil {
		s.Logger.Error("Error adding mod")
		s.Logger.Error(err.Error())
		return nil, err
	}

	return mapAddArchiveResponse(result), nil
}

func (s *ServiceHandlerServer) GetAllMods(ctx context.Context,
	req *servicev1.GetAllModsRequest,
) (*servicev1.GetAllModsResponse, error) {
	result, err := s.Handler.GetAllMods()
	if err != nil {
		s.Logger.Error("Failed to get all mods",
			"error", err,
		)
		return nil, err
	}

	s.Logger.Info("Mods fetched",
		"count", len(result.Mods),
	)
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
		s.Logger.Error("failed to map installed mods",
			"error", err,
		)
		return nil, err
	}
	s.Logger.Info("Preparing mod install",
		"modId", modID,
		"modVersionUUID", modVersionUUID,
		"choices.count", len(choices),
	)
	if err := s.Handler.InstallMod(modID, modVersionUUID, choices); err != nil {
		s.Logger.Error("failed to install mod",
			"err", err,
		)
		return nil, err
	}

	return &servicev1.InstallModResponse{}, nil
}

func (s *ServiceHandlerServer) UninstallMod(ctx context.Context,
	req *servicev1.UninstallModRequest,
) (*servicev1.UninstallModResponse, error) {
	ids := req.GetModIds()
	modIds := make([]int, len(ids))

	for _, mid := range ids {
		modIds = append(modIds, int(mid))
	}

	s.Logger.Info("Attempting to uninstall Mods",
		"ids.count", len(ids),
		"ids", modIds,
	)
	if err := s.Handler.UninstallMod(modIds); err != nil {
		s.Logger.Error("Failed to uninstall mods",
			"error", err,
		)
		return nil, err
	}

	return &servicev1.UninstallModResponse{}, nil
}

func (s *ServiceHandlerServer) DeleteMod(ctx context.Context,
	req *servicev1.DeleteModRequest,
) (*servicev1.DeleteModResponse, error) {
	toDeleteLen := len(req.GetModsToDelete())
	modsToDelete := make([]actions.DeleteModEntry, toDeleteLen)

	for _, td := range req.GetModsToDelete() {
		modsToDelete = append(modsToDelete, actions.DeleteModEntry{
			ID:             int(td.GetId()),
			DeleteVersions: mapDeleteModType(td.GetDeleteStates()),
			DeleteArchive:  td.GetDeleteArchive(),
		})
	}
	s.Logger.Info("Attempting to uninstall mods",
		"count", toDeleteLen,
	)
	if err := s.Handler.DeleteMod(actions.DeleteModRequest{Mods: modsToDelete}); err != nil {
		s.Logger.Error("Failed to delete mods",
			"error", err,
		)
		return nil, err
	}

	s.Logger.Info("Mods deleted")
	return &servicev1.DeleteModResponse{}, nil
}
