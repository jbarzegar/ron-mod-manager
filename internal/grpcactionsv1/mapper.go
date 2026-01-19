package grpcactionsv1

import (
	"github.com/google/uuid"
	"github.com/jbarzegar/ron-mod-manager/ent/mod"
	"github.com/jbarzegar/ron-mod-manager/handler"
	"github.com/jbarzegar/ron-mod-manager/internal/actions"
	servicev1 "github.com/jbarzegar/ron-mod-manager/proto/service/v1"
)

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
		ModId:        int32(result.Mod.ID),
		ModVersionId: result.ModVersion.UUID.String(),
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

// mapGetAllModsResponse map response to grpc readable
func mapGetAllModsResponse(result *actions.AllModsResponse) servicev1.GetAllModsResponse {
	data := make([]*servicev1.Mod, len(result.Mods))
	for _, mod := range result.Mods {
		data = append(data,
			&servicev1.Mod{
				Id:            int32(mod.ID),
				Name:          mod.Name,
				State:         mapStateToEnum(mod.State),
				ActiveVersion: &mod.ActiveVersion,
				Origin:        mod.Origin,
			},
		)
	}
	return servicev1.GetAllModsResponse{
		Data: data,
	}
}

func mapInstallModRequest(req *servicev1.InstallModRequest) (
	modID int, modVersionUUID uuid.UUID, choices []uuid.UUID, err error,
) {
	modID = int(req.GetModId())
	modVersionUUID, err = uuid.Parse(req.ModVersionId)
	if err != nil {
		return -1, uuid.NullUUID{}.UUID, []uuid.UUID{}, err
	}

	choices = make([]uuid.UUID, len(req.GetChoicesToActivate()))
	for _, id := range req.GetChoicesToActivate() {
		u, err := uuid.Parse(id)
		if err != nil {
			return -1, uuid.NullUUID{}.UUID, []uuid.UUID{}, err
		}
		choices = append(choices, u)
	}

	return modID, modVersionUUID, choices, nil
}

func mapDeleteModType(s servicev1.DeleteModType) string {
	switch s {
	case servicev1.DeleteModType_DELETE_MOD_TYPE_ACTIVE:
		return "active"
	case servicev1.DeleteModType_DELETE_MOD_TYPE_ALL:
		return "all"
	case servicev1.DeleteModType_DELETE_MOD_TYPE_UNSPECIFIED:
		return "all"
	default:
		return "all"
	}
}
