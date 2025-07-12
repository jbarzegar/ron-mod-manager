package handlerio

import (
	"errors"
	"io/fs"
	"path/filepath"

	"github.com/jbarzegar/ron-mod-manager/archive"
)

type FileSystemHandler struct{}

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

func (h *FileSystemHandler) InstallMod(archivePath string, outPath string) error {

	return errors.New("Not implemented")
}
