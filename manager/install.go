package manager

import (
	"fmt"
	"log"
	"path"
	"path/filepath"
	"strings"

	"github.com/jbarzegar/ron-mod-manager/config"
	"github.com/jbarzegar/ron-mod-manager/paths"
	statemanagement "github.com/jbarzegar/ron-mod-manager/state-management"
	"github.com/jbarzegar/ron-mod-manager/types"
	"github.com/jbarzegar/ron-mod-manager/utils"
)

func Install(n string) {
	state := statemanagement.GetState()
	absArchivePath := path.Join(paths.AbsArchiveDir(), n)

	archive, idx, _ := statemanagement.GetArchiveByName(n)

	absModPath := path.Join(paths.AbsModsDir(), archive.Name)

	err := utils.ExtractArchive(absArchivePath, absModPath, false)
	if err != nil {
		log.Fatal(err)
	}

	// Extract paks
	// List all paks in file
	// TODO: Make this recursive, rn it will only do a shallow check
	matches, _ := filepath.Glob(path.Join(absModPath, "*.pak"))

	mod := types.ModInstall{ArchiveName: n, Name: archive.Name, Paks: []string{}, State: "inactive"}

	for _, m := range matches {
		mod.Paks = append(mod.Paks, strings.Split(strings.Split(m, archive.Name)[1], "/")[1])
	}

	// Look through state to see if mod is already installed
	installed := false
	for _, x := range state.Mods {
		if x.Name == archive.Name {
			archive.Installed = true
			installed = true
			break
		}
	}

	if !installed {
		state.Archives[idx].Installed = true
		// append mod and write to state
		state.Mods = append(state.Mods, mod)
		statemanagement.WriteState(state, config.GetConfig())

		fmt.Println("Mod installed")
	} else {
		fmt.Println("Mod already installed")
	}

}
