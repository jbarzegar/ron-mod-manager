// handler contains the user specific methods needed for interfacing with the mod manager
// contains logic for standard mod manager activites
package handler

import (
	"context"
	"errors"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jbarzegar/ron-mod-manager/appconfig"
	"github.com/jbarzegar/ron-mod-manager/archive"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/mod"
	"github.com/jbarzegar/ron-mod-manager/handlerio"
)

type AddModResponse struct {
	// *ent.Archive
	Archive    *ent.Archive    `json:"archive"`
	ModVersion *ent.ModVersion `json:"modVersion"`
	// a list of potential paks that can be added
	Choices []archive.Choice `json:"choices"`
}

type handler interface {
	// get all mods avaiable to the db
	GetMods()
	// add an archive as a mod
	AddMod(archivePath string) (AddModResponse, error)
	// install a mod
	InstallMod()
}

// defines a Handler struct
// a Handler struct takes in a number of dependencies
// used for data persistence and working with IO
type Handler struct {
	// instance of a ent db client
	db *ent.Client
	// instance of the application config
	config *appconfig.AppConfig
	// Defines the logic handling IO
	// IO can mean different things in Different contexts
	// Generally in regular usecases IO refers to FileSystemHandler
	// abstracting this enables this portion of code to be more easily testable
	io handlerio.IOHandler
}

// implement handler with dep inversion
func NewHandler(db *ent.Client, config *appconfig.AppConfig, ioHandler handlerio.IOHandler) Handler {
	h := Handler{db: db, config: config, io: ioHandler}
	return h
}

func createModVersionBuilder(
	db *ent.Client,
	name string,
	version string,
	archives ent.Archives,
) (*ent.ModVersionCreate, uuid.UUID, error) {
	uuid, err := uuid.NewV6()
	if err != nil {
		return nil, uuid, err
	}
	// create a new mod version
	modVersionBuilder := db.ModVersion.Create().
		SetName(name).
		SetVersion(version).
		SetUUID(uuid).AddArchives(archives...)

	return modVersionBuilder, uuid, nil
}

type modBuilderPayload struct {
	Name    string
	Version string
	Origin  string
}

func createModBuilder(
	_ context.Context,
	db *ent.Client,
	m modBuilderPayload,
	modVersionUUID uuid.UUID,
) (*ent.ModCreate, error) {
	// create a new mod the new mod should be inactive since no instance of
	// this mod exists we can safely(?probablynotsafe) assign the current mod version to the
	// above uuid
	modBuilder := db.Mod.Create().
		SetName(m.Name).
		SetState("inactive").
		SetActiveVersion(modVersionUUID).
		SetOrigin(m.Origin)

	return modBuilder, nil
}

// createPakBuilders will map a slice of choices, create an ent pakBuilder for each choice
// returning a bulk create builder
func createPakBuilders(db *ent.Client, choices []archive.Choice) (*ent.PakCreateBulk, error) {
	// create slice of pak builders
	builders := []*ent.PakCreate{}
	for _, c := range choices {
		uid, err := uuid.NewV6()
		if err != nil {
			return nil, err
		}
		pakBuilder := db.Pak.Create().
			SetUUID(uid).
			SetActive(false).
			SetName(c.Name).
			SetPath(c.FullPath)

		builders = append(builders, pakBuilder)
	}

	return db.Pak.CreateBulk(builders...), nil
}

// AddMod adds a mod to the current mod manager registry
func (h *Handler) AddMod(archivePath string, name string) (*AddModResponse, error) {
	p := filepath.Join(h.config.ModDir, "archives", name)
	ctx := context.TODO()
	archive, err := h.db.Archive.Create().
		SetName(name).
		SetArchivePath(p).
		// todo calc md5 sum
		SetMd5Sum("").
		Save(ctx)
	if err != nil {
		return nil, errors.Join(
			errors.New("Error creating archive"),
			err,
		)
	}

	// handle IO portion of addmod after db operations have finished
	choices, err := h.io.AddMod(archivePath, p)
	if err != nil {
		return nil, errors.Join(
			errors.New("IO Error while adding mod"),
			err,
		)
	}

	// create and save paks to db
	// there's no way of knowing what paks are avaliable until you look at them from the IO layer
	paksBuilder, err := createPakBuilders(h.db, choices)
	if err != nil {
		return nil, errors.Join(errors.New("error while building PAKs"), err)
	}

	// save the paks
	paks, err := paksBuilder.Save(context.TODO())
	if err != nil {
		return nil, errors.Join(errors.New("error saving paks"), err)
	}

	modVersionBuilder, uid, err := createModVersionBuilder(
		h.db,
		name,
		"0.0.0",
		ent.Archives{archive},
	)
	if err != nil {
		return nil, errors.Join(
			errors.New("error while creating modVersionBuilder"),
			err,
		)
	}

	//append paks to the mod version builder
	modVersionBuilder.AddPaks(paks...)

	// save the builder
	modVersion, err := modVersionBuilder.Save(context.TODO())
	if err != nil {
		return nil, errors.Join(errors.New("error saving mod version"), err)
	}

	modExists, err := h.db.Mod.Query().Where(mod.NameEQ(name)).Exist(ctx)
	if err != nil {
		return nil, errors.Join(
			errors.New("error getting mod"),
			err,
		)
	}
	// we need to create a new mod entry if it doesn't exist already
	if !modExists {
		newModBuilder, err := createModBuilder(
			ctx,
			h.db,
			modBuilderPayload{
				Name:    name,
				Version: "0.0.0",
				Origin:  "unknown",
			},
			uid,
		)
		if err != nil {
			return nil, errors.Join(errors.New("error creating modBuilder"), err)
		}

		newModBuilder.AddVersionIDs(modVersion.ID)

		// save the new mod
		_, err = newModBuilder.Save(ctx)
		if err != nil {
			return nil, errors.Join(errors.New("error saving mod"), err)
		}
	}

	return &AddModResponse{
		Archive:    archive,
		ModVersion: modVersion,
		Choices:    choices,
	}, nil
}

func (h *Handler) InstallMod() error {
	// h.db.Mod.
	//
	// re
	return nil
}
