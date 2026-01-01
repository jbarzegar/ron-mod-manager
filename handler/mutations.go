package handler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/jbarzegar/ron-mod-manager/archive"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/mod"
	"github.com/jbarzegar/ron-mod-manager/ent/modversion"
	"github.com/jbarzegar/ron-mod-manager/ent/pak"
	"github.com/jbarzegar/ron-mod-manager/handlerio"
)

func createModVersionBuilder(
	db *ent.Client,
	name string,
	version string,
	archives ent.Archives,
) (*ent.ModVersionCreate, uuid.UUID, error) {
	uuid := uuid.New()
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
		SetState(mod.StateInactive).
		SetActiveVersion(modVersionUUID).
		SetOrigin(m.Origin)

	return modBuilder, nil
}

// createPakBuilders will map a slice of choices, create an ent pakBuilder for each choice
// returning a bulk create builder
func createPakBuilders(db *ent.Client, choices []archive.Choice) (*ent.PakCreateBulk, error) {
	// create slice of pak builders
	builders := []*ent.PakCreate{}
	for i, c := range choices {
		uid := uuid.New()

		// TODO: Remodel. this feels gross
		// because it is
		choices[i].UUID = uid
		pakBuilder := db.Pak.Create().
			SetUUID(uid).
			SetActive(false).
			SetName(c.Name).
			SetPath(c.FullPath)

		builders = append(builders, pakBuilder)
	}

	return db.Pak.CreateBulk(builders...), nil
}

// AddArchive adds a mod to the current mod manager registry
func (h *Handler) AddArchive(archivePath string, name string) (*AddModResponse, error) {
	version := "0.0.0"
	p := filepath.Join(h.Config.ModDir, "mods", name, version)
	ctx := context.TODO()
	archive, err := h.Db.Archive.Create().
		SetName(name).
		SetArchivePath(archivePath).
		// todo calc md5 sum
		SetMd5Sum("").
		Save(ctx)
	if err != nil {
		if ent.IsConstraintError(err) {
			slog.Warn("archive already exists. TODO: ATTEMPT TO INSTALL ANYWAY")
			return nil, fmt.Errorf("Archive: %v already exists as an archive", name)
		}
		return nil, errors.Join(
			errors.New("Error creating archive"),
			err,
		)
	}

	// handle IO portion of addmod after db operations have finished
	choices, err := h.Io.AddArchive(archivePath, p)
	if err != nil {
		return nil, errors.Join(
			errors.New("IO Error while adding mod"),
			err,
		)
	}

	// create and save paks to db
	// there's no way of knowing what paks are avaliable until you look at them from the IO layer
	paksBuilder, err := createPakBuilders(h.Db, choices)
	if err != nil {
		return nil, errors.Join(errors.New("error while building PAKs"), err)
	}

	// save the paks
	paks, err := paksBuilder.Save(context.TODO())
	if err != nil {
		return nil, errors.Join(errors.New("error saving paks"), err)
	}

	modVersionBuilder, uid, err := createModVersionBuilder(
		h.Db,
		name,
		version,
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

	modExists, err := h.Db.Mod.Query().Where(mod.NameEQ(name)).Exist(ctx)
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
			h.Db,
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
		mod, err := newModBuilder.Save(ctx)
		if err != nil {
			return nil, errors.Join(errors.New("error saving mod"), err)
		}

		return &AddModResponse{
			Mod:        mod,
			Archive:    archive,
			ModVersion: modVersion,
			Choices:    choices,
		}, nil
	}

	return nil, errors.New("mod exists already. what do?")

}

// InstallMod activates a given modVersion to a given mod including paks that
// must be activated as well
func (h *Handler) InstallMod(modID int, modVersionUUID uuid.UUID, paksToActivate []uuid.UUID) error {
	// get m we want to update the modversion on
	m, err := h.Db.Mod.Query().Where(mod.ID(modID)).Only(context.TODO())
	if err != nil {
		return errors.Join(
			fmt.Errorf("failed to get mod with id %v", modID),
			err,
		)
	}

	mv, err := h.
		Db.
		ModVersion.
		Query().
		Where(modversion.UUIDEQ(modVersionUUID)).
		Only(context.TODO())
	if err != nil {
		return errors.Join(
			fmt.Errorf("could not get mod version %v", modVersionUUID),
			err,
		)
	}

	paks := []*ent.Pak{}
	// set the paks that must be activated
	for _, pUUID := range paksToActivate {
		Pak := h.Db.Pak
		pakEQ := pak.UUIDEQ(pUUID)
		exists, err := Pak.Query().Where(pakEQ).Exist(context.TODO())
		if err != nil {
			return errors.Join(
				errors.New("failed to update pak"),
				err,
			)
		}
		if exists {
			pak, err := Pak.Query().Where(pakEQ).Only(context.Background())
			if err != nil {
				return errors.Join(
					fmt.Errorf("Error querying updated pak: %v", pUUID),
					err,
				)
			}
			paks = append(paks, pak)
		}
	}

	modUpdate := m.Update()

	// set the active version of the modVersion to the mod
	modUpdate.SetActiveVersion(mv.UUID)

	// set the state to active if it's not already set
	modUpdate.SetState(mod.StateActivate)

	// finally save the mod
	_, err = modUpdate.Save(context.TODO())
	if err != nil {
		return errors.Join(
			errors.New("failed updating mod"),
			err,
		)
	}

	// do the io work
	if err := h.Io.InstallMod(handlerio.Installable{
		Mod: m,
		Assets: handlerio.InstallableAssets{
			Paks: paks,
		},
		OutPath: h.Config.StagingModFolderName,
	}); err != nil {
		return err
	}

	return nil
}
