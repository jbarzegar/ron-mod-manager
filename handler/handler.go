// handler contains the user specific methods needed for interfacing with the mod manager
// contains logic for standard mod manager activites
package handler

import (
	"context"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jbarzegar/ron-mod-manager/appconfig"
	"github.com/jbarzegar/ron-mod-manager/archive"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/mod"
	"github.com/jbarzegar/ron-mod-manager/ent/modversion"
	"github.com/jbarzegar/ron-mod-manager/ent/pak"
	"github.com/jbarzegar/ron-mod-manager/handlerio"
	"github.com/jbarzegar/ron-mod-manager/internal/actions"
)

type AddModResponse struct {
	Mod        *ent.Mod        `json:"mod"`
	Archive    *ent.Archive    `json:"archive"`
	ModVersion *ent.ModVersion `json:"modVersion"`
	// a list of potential paks that can be added
	Choices []archive.Choice `json:"choices"`
}

type handler interface {
	// get all mods avaiable to the db
	GetAllMods() (actions.AllModsResponse, error)
	// get all mods currently marked as "staged"
	// staged mods are ones being loaded into the game at this moment
	GetStagedMods() (actions.StagedResponse, error)
	// add an archive as a mod
	AddArchive(archivePath string) (AddModResponse, error)
	AddModByID(modId int) (AddModResponse, error)
	// install a mod using an instance of a mod version
	InstallMod(modID int, modVersion *ent.ModVersion, paksToActivate []uuid.UUID) error
	// uninstall a mod from the staging area
	UninstallMod(modID []int) error
}

// defines a Handler struct
// a Handler struct takes in a number of dependencies
// used for data persistence and working with IO
type Handler struct {
	// instance of a ent Db client
	Db *ent.Client
	// instance of the application Config
	Config *appconfig.AppConfig
	// Defines the logic handling IO
	// IO can mean different things in Different contexts
	// Generally in regular usecases IO refers to FileSystemHandler
	// abstracting this enables this portion of code to be more easily testable
	Io handlerio.IOHandler
}

// implement handler with dep inversion
func NewHandler(db *ent.Client, config *appconfig.AppConfig, ioHandler *handlerio.IOHandler) Handler {
	// io := *ioHandler
	h := Handler{Db: db, Config: config, Io: *ioHandler}
	return h
}

func (h *Handler) GetAllMods() (*actions.AllModsResponse, error) {
	allMods, err := h.Db.Mod.Query().All(context.TODO())
	if err != nil {
		return nil, err
	}

	resp := &actions.AllModsResponse{}
	for _, m := range allMods {
		versions, err := h.Db.ModVersion.Query().All(context.TODO())
		if err != nil {
			return nil, err
		}

		mappedMod := actions.AllModsEntry{
			Name:          m.Name,
			State:         m.State.String(),
			Origin:        m.Origin,
			ActiveVersion: m.ActiveVersion.String(),
			Versions:      versions,
		}

		resp.Mods = append(resp.Mods, mappedMod)
	}

	return resp, nil
}

func (h *Handler) GetArchives(req *actions.GetArchiveRequest) (*actions.GetArchivesResponse, error) {
	resp := &actions.GetArchivesResponse{}
	archives, err := h.Db.Archive.Query().All(context.TODO())
	if err != nil {
		return nil, err
	}

	registeredPaths := []string{}
	for _, archive := range archives {
		a := actions.GetArchivesEntry{
			ID:                 archive.ID,
			Name:               archive.Name,
			Path:               archive.ArchivePath,
			Installed:          archive.Installed,
			PathExists:         false,
			ValidationMessages: []actions.ValidationMessage{},
		}

		// Path exists
		a.PathExists = h.Io.PathExists(a.Path)

		// TODO: create validation messages
		// -- mod installed but path doesn't exist
		if a.Installed && !a.PathExists {
			a.ValidationMessages = append(
				a.ValidationMessages,
				actions.ValidationMessage{
					Kind:    "warning",
					Message: "Archive installed, but path does not exist in path anymore",
				},
			)
		}
		// -- neither installed nor has path that exists
		if !a.Installed && !a.PathExists {
			a.ValidationMessages = append(a.ValidationMessages,
				actions.ValidationMessage{
					Kind:    "info",
					Message: "Archive neither installed nor the archive path found. ",
				},
			)
		}

		a.Installable = a.PathExists

		// --
		registeredPaths = append(registeredPaths, a.Path)
		resp.Archives = append(resp.Archives, a)
	}

	if req.Untracked {
		dirs, err := h.Io.GetUntracked(registeredPaths)
		if err != nil {
			return nil, err
		}

		for _, d := range dirs {
			resp.UntrackedArchives = append(resp.UntrackedArchives, actions.GetArchivesEntry{
				ID:   -1,
				Name: d.Name,
				Path: d.Path,
				// we can assume this exists. Since this is read directly from the file system
				// the os itself assumes we can install it
				Installable: true,
				Installed:   false,
				PathExists:  true,
			})

		}
	}

	return resp, nil
}

func (h *Handler) GetStagedMods() (*actions.StagedResponse, error) {
	// get mods marked as active
	activeMods, err := h.Db.Mod.Query().
		Where(mod.StateEQ(mod.StateActivate)).All(context.TODO())
	if err != nil {
		return nil, err
	}

	resp := &actions.StagedResponse{}
	for _, m := range activeMods {
		av, err := h.Db.ModVersion.Query().
			Where(modversion.UUIDEQ(*m.ActiveVersion)).Only(context.Background())
		if err != nil {
			return nil, err
		}

		activePaks, err := av.QueryPaks().Where(pak.ActiveEQ(true)).All(context.Background())
		if err != nil {
			return nil, err
		}
		mappedMod := actions.StagedMod{
			Name:              m.Name,
			State:             string(m.State),
			Origin:            m.Origin,
			ActiveVersion:     av.Version,
			ActiveVersionUUID: m.ActiveVersion,
			Paks:              activePaks,
		}

		resp.Mods = append(resp.Mods, mappedMod)
	}

	return resp, nil
}

func (h *Handler) UninstallMod(modIDs []int) error {
	for _, id := range modIDs {
		// get the mod
		m, err := h.Db.Mod.Query().Where(mod.IDEQ(id), mod.StateEQ(mod.StateActivate)).Only(context.Background())
		if err != nil {
			return err
		}

		modBuilder := m.Update().SetState(mod.StateInactive)

		// query the active vesrion of the mod & grab active paks
		paks, err := m.QueryVersions().
			Where(modversion.UUIDEQ(*m.ActiveVersion)).
			QueryPaks().
			Where(pak.ActiveEQ(true)).All(context.TODO())
		if err != nil {
			return err
		}

		if _, err := modBuilder.Save(context.TODO()); err != nil {
			return err
		}

		if err := h.Io.UninstallMod(paks); err != nil {
			return err
		}
	}

	return nil
}

func (h *Handler) DeleteMod(req actions.DeleteModRequest) error {
	db := h.Db
	for _, m := range req.Mods {
		md, err := db.Mod.Get(context.TODO(), m.ID)
		if err != nil {
			return err
		}

		// then delete everything
		mvs, err := md.QueryVersions().All(context.TODO())
		if err != nil {
			return err
		}

		// operate on mod versions
		// TODO: allow for selectively deleting some versions instead of nuking everything by default
		for _, mv := range mvs {
			// query the paks
			ids, err := mv.QueryPaks().All(context.TODO())
			if err != nil {
				return err
			}
			// unstage the mod if it's active
			if md.State == mod.StateActivate {
				h.Io.UninstallMod(ids)
			}
			for _, id := range ids {
				if _, err := db.Pak.Delete().Where(pak.IDEQ(id.ID)).Exec(context.TODO()); err != nil {
					return err
				}

			}

			// delete the mod versions
			if _, err := db.ModVersion.Delete().Where(modversion.IDEQ(mv.ID)).Exec(context.TODO()); err != nil {
				return err
			}
		}

		// delete the actual mod now
		if _, err := db.Mod.Delete().Where(mod.IDEQ(m.ID)).Exec(context.TODO()); err != nil {
			return err
		}

		// deleting the mod on the io level is the last step
		if err := h.Io.DeleteMod(filepath.Join(h.Config.ModDir, "mods", md.Name)); err != nil {
			return err
		}
	}
	return nil
}
