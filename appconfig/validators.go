package appconfig

import (
	"errors"
	"fmt"
	"os"
	"path"
)

func validateRonDir(dir string) error {
	if dir == "unknown" {
		return errors.New("dir not set")
	}

	dirsToCheck := [3]string{"ReadyOrNot.exe", "Engine", "ReadyOrNot"}
	invalid := []string{}

	for _, x := range dirsToCheck {
		if _, err := os.Stat(path.Join(dir, x)); os.IsNotExist(err) {
			invalid = append(invalid, dir)
		}
	}

	if len(invalid) > 0 {
		return errors.New("invalid RON dir")
	}

	return nil
}

func validateModDir(dir string) error {
	// Archives stores the .zips used for installs
	archivesPath := path.Join(dir, "archives")
	// Stored "installed mods"
	modsPath := path.Join(dir, "mods")

	dirsToEnsure := []string{archivesPath, modsPath}
	for _, d := range dirsToEnsure {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			fmt.Println(d, "doesn't exist, creating")
			err = os.MkdirAll(d, 0700)
			if err != nil {
				return err
			}

		}
	}

	return nil
}
