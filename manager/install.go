package manager

import (
	"fmt"
	"log"
	"os"
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

	var matches []string

	filepath.Walk(absModPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			if err != nil {
				log.Fatal(err)
				return err
			}

			if strings.HasSuffix(path, ".pak") {
				matches = append(matches, path)
			}

		}

		return nil
	})

	mod := types.ModInstall{ArchiveName: n, Name: archive.Name, Paks: []string{}, State: "inactive"}

	for _, m := range matches {
		relPakPath := strings.Split(m, absModPath+"/")[1]
		mod.Paks = append(mod.Paks, relPakPath)
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
