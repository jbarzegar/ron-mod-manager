package statemanagement

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/iancoleman/strcase"
	"github.com/jbarzegar/ron-mod-manager/config"
	"github.com/jbarzegar/ron-mod-manager/db"
	"github.com/jbarzegar/ron-mod-manager/ent/archive"
	"github.com/jbarzegar/ron-mod-manager/types"
	"github.com/jbarzegar/ron-mod-manager/utils"
)

func writeConfigFile(confFile string, config types.MMConfig) {
	b, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}

	err = os.WriteFile(confFile, b, 0666)

	if err != nil {
		log.Fatal(err)
	}
}

func validateRonDir(d string) error {
	dirsToCheck := [3]string{"ReadyOrNot.exe", "Engine", "ReadyOrNot"}
	invalid := []string{}

	for _, x := range dirsToCheck {
		if _, err := os.Stat(path.Join(d, x)); os.IsNotExist(err) {
			invalid = append(invalid, d)

		}
	}

	if len(invalid) > 0 {
		return errors.New("invalid RON dir")
	}

	return nil
}

func validateModDir(c string) {
	// Archives stores the .zips used for installs
	archivesPath := path.Join(c, "archives")
	// Stored "installed mods"
	modsPath := path.Join(c, "mods")
	//
	// modState := path.Join(c, "mod-state.meta.json")

	dirsToEnsure := []string{archivesPath, modsPath}
	for _, d := range dirsToEnsure {
		if _, err := os.Stat(d); os.IsNotExist(err) {
			fmt.Println(d, "doesn't exist, creating")
			err = os.MkdirAll(d, 0700)

			if err != nil {
				log.Fatal(err)
			}

		}
	}
}

func ensureConfig(confPath string) types.MMConfig {
	defaultConfig := types.MMConfig{GameDir: "unknown", ModDir: "unknown"}

	if _, err := os.Stat(confPath); errors.Is(err, os.ErrNotExist) {
		fmt.Println("Write config file")
		writeConfigFile(confPath, defaultConfig)
	} else {
		c := config.ReadConfFile(confPath)

		err := validateRonDir(c.GameDir)

		if err != nil {
			log.Fatal(err)
		}
		validateModDir(c.ModDir)

		return c
	}
	return defaultConfig

}

func genMd5Sums(p any) string {
	buf := &bytes.Buffer{}

	gob.NewEncoder(buf).Encode(p)

	sum := md5.Sum(buf.Bytes())

	hash := fmt.Sprintf("%x", sum)

	return hash
}

const stateFile = "mod-state.meta.json"

// Load _state into memory
var _state types.ModState

func GetState() types.ModState {
	conf := config.GetConfig()

	file, _ := os.ReadFile(path.Join(conf.ModDir, stateFile))

	json.Unmarshal(file, &_state)

	return _state

}

func WriteState(s types.ModState, c types.MMConfig) {
	sf := path.Join(c.ModDir, stateFile)
	json, _ := json.Marshal(s)

	os.WriteFile(sf, json, 0666)

}

func listArchives(modDir string) []string {
	dirs := []string{"zip", "rar", "7z"}
	var archives []string

	for _, ext := range dirs {
		g := path.Join(modDir, "archives", "*."+ext)

		d, _ := filepath.Glob(g)

		archives = append(archives, d...)
	}

	return archives
}

// cecks or state setup prior to running the application
func PreflightChecks() {
	// var ex string
	ex := config.ConfPath

	// Detect and setup config
	// Creates dir structure
	// Generates initial mod state file
	configFilePath := path.Join(ex, "ron-mm.conf.json")

	err := os.MkdirAll(ex, 0700)

	if err != nil {
		log.Fatalf("Err creating config path", err)
	}

	c := ensureConfig(configFilePath)

	archives := listArchives(c.ModDir)
	archiveModel := db.Client().Archive

	for _, archivePath := range archives {
		name := strcase.ToKebab(utils.SplitArchivePath(utils.FormatArchiveName(archivePath)))

		a := archiveModel.
			Query().Where(archive.Name(name))

		archiveExists, err := a.Exist(context.Background())

		if err != nil {
			log.Fatalf("Failure reading archive", err)
		}

		if !archiveExists {
			// Write to disk
			fmt.Println(archivePath)
			_, err := archiveModel.
				Create().
				SetName(name).
				SetArchivePath(archivePath).
				SetMd5Sum(genMd5Sums(name)).
				SetInstalled(false).
				Save(context.Background())

			if err != nil {
				log.Fatalf("Failure writing archive", err)
			}
		}
	}
}

func GetArchiveByName(name string) (*types.Archive, int, error) {
	state := GetState()
	var selectedArchive *types.Archive = nil
	var selectedArchiveIdx = -1
	for i, m := range state.Archives {
		if m.ArchiveFile == name {
			selectedArchive = &m
			selectedArchiveIdx = i
			break
		}
	}

	return selectedArchive, selectedArchiveIdx, nil
}

func GetModByName(name string) (*types.ModInstall, int, error) {
	state := GetState()

	var selectedMod *types.ModInstall = nil
	var idx = -1
	for i, m := range state.Mods {
		if m.Name == name {
			selectedMod = &m
			idx = i
			break
		}
	}

	if selectedMod == nil {
		return nil, -1, errors.New("mod not found")
	}

	return selectedMod, idx, nil
}

// Get
func GetModsByState(filter string) []types.ModInstall {
	state := GetState()

	var choices []types.ModInstall
	for _, mod := range state.Mods {
		switch filter {
		case "active", "inactive":
			if mod.State == filter {
				choices = append(choices, mod)
			}
		case "":
			choices = append(choices, mod)

		// handle unknown filters
		default:
			fmt.Println("WARN unsupported filter:", filter, " handling as if unfiltered")
			choices = append(choices, mod)

		}
	}

	return choices

}
