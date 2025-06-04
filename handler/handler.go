// handler contains the user specific methods needed for interfacing with the mod manager
// contains logic for standard mod manager activites
package handler

import (
	"context"
	"path/filepath"

	"github.com/jbarzegar/ron-mod-manager/appconfig"
	"github.com/jbarzegar/ron-mod-manager/archive"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/handlerio"
)

type AddModResponse struct {
	// *ent.Archive
	Archive *ent.Archive `json:"archive"`
	// a list of potential paks that can be added
	Choices []archive.Choice `json:"choices"`
}

type handler interface {
	// get all mods avaiable to the db
	GetMods()
	// add an archive as a mod
	AddMod(archivePath string) (AddModResponse, error)
	InstallMod()
}

// defines a Handler struct
// a Handler struct takes in a number of dependencies
// used for data persistence and working with IO
type Handler struct {
	db     *ent.Client
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

// AddMod adds a mod to the current mod manager registry
func (h *Handler) AddMod(archivePath string, name string) (*AddModResponse, error) {
	p := filepath.Join(h.config.ModDir, "archives", name)

	ctx := context.TODO()
	a, err := h.db.Archive.Create().
		SetName(name).
		SetArchivePath(p).
		// todo calc md5 sum
		SetMd5Sum("").
		Save(ctx)
	if err != nil {
		return nil, err
	}

	// handle IO portion of addmod
	choices, err := h.io.AddMod(archivePath, p)
	if err != nil {
		return nil, err
	}

	return &AddModResponse{
		Archive: a,
		Choices: choices,
	}, nil
}
