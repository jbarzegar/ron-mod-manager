package handlerio

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/jbarzegar/ron-mod-manager/appconfig"
	"github.com/jbarzegar/ron-mod-manager/archive"
	"github.com/jbarzegar/ron-mod-manager/ent"
)

type FileSystemHandler struct {
	Config *appconfig.AppConfig
}

// collectChoices gets all possible choices from an archive
func collectChoices(p string) []archive.Choice {
	choices := []archive.Choice{}
	filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
		if !d.Type().IsDir() {
			ext := filepath.Ext(path)
			_, f := filepath.Split(path)

			if ext == ".pak" {
				choices = append(choices, archive.Choice{
					Name:     f,
					FullPath: path,
				})
			}
		}
		return err
	})

	return choices
}

func (h *FileSystemHandler) AddMod(archivePath string, outputPath string) ([]archive.Choice, error) {
	if err := archive.Extract(archivePath, outputPath, true); err != nil {
		return nil, err
	}

	// get choices
	choices := collectChoices(outputPath)

	return choices, nil
}

func (h *FileSystemHandler) InstallMod(payload Installable) error {
	// install mods to the staging folder
	for _, x := range payload.Assets.Paks {
		pakPath := filepath.Join(h.Config.ModDir, h.Config.StagingModFolderName, x.Name)
		f, err := os.Open(x.Path)
		if err != nil {
			return errors.Join(
				fmt.Errorf("failed to open file %v", x.Path),
				err,
			)
		}
		defer f.Close()
		// copy the paks around
		w, err := os.Create(pakPath)
		if err != nil {
			return errors.Join(
				fmt.Errorf("failed to create staged file %v", pakPath),
				err,
			)
		}
		// install the paks
		_, err = io.Copy(w, f)
		if err != nil {
			return errors.Join(
				fmt.Errorf("failed to copy staged file %v", pakPath),
				err,
			)
		}
	}
	return nil
}

func (h *FileSystemHandler) UninstallMod(pakPaths []*ent.Pak) error {
	for _, x := range pakPaths {
		pakPath := filepath.Join(h.Config.ModDir, h.Config.StagingModFolderName, x.Name)
		if err := os.Remove(pakPath); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			} else {
				return err
			}
		}
	}

	return nil
}

func (h *FileSystemHandler) DeleteMod(modPath string) error {
	return os.RemoveAll(modPath)
}
