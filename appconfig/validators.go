package appconfig

import (
	"errors"
	"fmt"
	"os"
	"path"
)

func validateRonDir(cfg AppConfig) error {
	dir := cfg.GameDir

	if dir == "unknown" {
		return errors.New("dir not set" + dir)
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

func validateModDir(cfg AppConfig) error {
	dir := cfg.ModDir
	modInstallDir := cfg.StagingModFolderName

	// Archives stores the .zips used for installs
	archivesPath := path.Join(dir, "archives")
	// Stored "installed mods"
	modsPath := path.Join(dir, "mods")
	// stores the staged mod files
	modInstallDirPath := path.Join(dir, modInstallDir)
	dirsToEnsure := []string{archivesPath, modsPath, modInstallDirPath}
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
