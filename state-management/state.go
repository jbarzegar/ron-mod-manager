package statemanagement

import (
	"bytes"
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
	"github.com/jbarzegar/ron-mod-manager/types"
	"github.com/jbarzegar/ron-mod-manager/utils"
)

func writeConfigFile(confFile string, config types.MMConfig) {
	b, err := json.Marshal(config)
	if err != nil {
		panic(err)
	}
	fmt.Println(b)

	os.WriteFile(confFile, b, 0666)
}

func validateRonDir(d string) error {
	dirsToCheck := [3]string{"ReadyOrNot.exe", "Engine", "ReadyOrNot"}
	invalid := []string{}

	for _, x := range dirsToCheck {
		fmt.Println(d, x)
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
	// filesToEnsure := []string{modState}

	for _, d := range dirsToEnsure {
		// if _, err := os.Stat(d)
		// x, _ := os.Stat(d)
		if _, err := os.Stat(d); os.IsNotExist(err) {
			fmt.Println(d, "doesn't exist")
		}

	}

	// for _, f := range filesToEnsure {
	// 	if _, err := os.Stat(f); os.IsNotExist(err) {
	// 		m, _ := json.Marshal(ModState{archiveSum: ""})
	// 		// os.WriteFile(f, m, 0666)
	// 	}
	// }

}

func ensureConfig(confPath string) types.MMConfig {
	defaultConfig := types.MMConfig{GameDir: "unknown", ModDir: "unknown"}

	if _, err := os.Stat(confPath); errors.Is(err, os.ErrNotExist) {
		// fmt.Println("Write config file")
		writeConfigFile(confPath, defaultConfig)
	} else {
		c := config.ReadConfFile(confPath)
		// var c types.MMConfig

		// f, err := os.ReadFile(confPath)
		// if err != nil {
		// 	panic(err)
		// }

		// json.Unmarshal(f, &c)

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

// Run through a large amount of checks or state setup prior to running the application
func PreflightChecks() {
	ex, err := os.Getwd()

	if err != nil {
		log.Fatal(err)
	}
	// Detect and setup config
	// Creates dir structure
	// Generates initial mod state file
	c := ensureConfig(path.Join(ex, "conf.json"))

	// Load state into memory
	s := GetState()

	// state.archiveSum

	// read archives dir generate md5 sum
	// add to mod-state.meta.json
	archives := listArchives(c.ModDir)

	// generate md5 sum, (skip if md5 matches one present in state file)
	sum := genMd5Sums(archives)

	if sum != s.ArchiveSum {
		// resync sum
		s.ArchiveSum = sum
		// resync archives and sums
		for _, archivePath := range archives {
			archiveName := utils.SplitArchivePath(archivePath)
			name := strcase.ToKebab(utils.SplitArchivePath(utils.FormatArchiveName(archivePath)))
			a := types.Archive{ArchiveFile: archiveName, Name: name, Md5Sum: genMd5Sums(archivePath), Installed: false}
			s.Archives = append(s.Archives, a)
		}
		// save state
		WriteState(s, c)
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

func GetModByName(name string) (*types.ModInstall, error) {
	state := GetState()

	var selectedMod *types.ModInstall = nil
	for _, m := range state.Mods {
		if m.Name == name {
			selectedMod = &m
			break
		}
	}

	if selectedMod == nil {
		return nil, errors.New("Mod not found")
	}

	return selectedMod, nil
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
