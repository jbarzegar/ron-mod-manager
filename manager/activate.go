package manager

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jbarzegar/ron-mod-manager/components"
	"github.com/jbarzegar/ron-mod-manager/config"
	"github.com/jbarzegar/ron-mod-manager/paths"
	statemanagement "github.com/jbarzegar/ron-mod-manager/state-management"
	"github.com/jbarzegar/ron-mod-manager/types"
)

func linkPak(m *types.ModInstall, p string, modName string) {

	state := statemanagement.GetState()
	absPakModPath := path.Join(paths.AbsModsDir(), m.Name, p)
	absPakGamePath := paths.PaksDir()

	_, err := os.Stat(path.Join(absPakGamePath, p))

	if !os.IsNotExist(err) {
		// Search all mods and see if the new mod was installed already
		// Account for mods that were installed by ron-mm
		for _, q := range state.Mods {
			if q.Name == modName && q.State == "active" {
				fmt.Println("Mod already installed")
				return
			}
		}
		// Mods may be installed manually and need to be accounted for
		fmt.Println("Mod installed outside of ron-mm's delete mod prior to activating")

	} else {
		for i, q := range state.Mods {
			if q.Name == modName {
				state.Mods[i].State = "active"
			}
		}

		statemanagement.WriteState(state, config.GetConfig())
		err = os.Symlink(absPakModPath, path.Join(absPakGamePath, p))

		if err != nil {
			log.Fatal("why", err)
		}

	}

}

func Activate(modsToActivate map[int]string) {
	fmt.Println("activating mods ")

	for _, modName := range modsToActivate {
		m, _ := statemanagement.GetModByName(modName)

		fmt.Println(m.Name, m.Paks)

		if len(m.Paks) == 1 {
			p := m.Paks[0]
			linkPak(m, p, modName)

		} else {
			choices := components.SelectMod(m.Paks)
			state := statemanagement.GetState()

			for _, p := range choices {
				// x := strings.Split(p, string(os.PathSeparator))
				// fmt.Println(, x[1])
				// linkPak(m, p, modName)

				absPakModPath := path.Join(paths.AbsModsDir(), m.Name, p)
				absPakGamePath := paths.PaksDir()

				x := path.Join(absPakGamePath, strings.Split(p, string(os.PathSeparator))[1])
				_, err := os.Stat(x)

				if !os.IsNotExist(err) {
					// Search all mods and see if the new mod was installed already
					// Account for mods that were installed by ron-mm
					for _, q := range state.Mods {
						if q.Name == modName && q.State == "active" {
							fmt.Println("Mod already installed")
							return
						}
					}
					// Mods may be installed manually and need to be accounted for
					fmt.Println("Mod installed outside of ron-mm's delete mod prior to activating")

				} else {
					for i, q := range state.Mods {
						if q.Name == modName {
							state.Mods[i].State = "active"
						}
					}

					statemanagement.WriteState(state, config.GetConfig())
					err = os.Symlink(absPakModPath, x)

					if err != nil {
						log.Fatal("why", err)
					}

				}

			}
		}

	}

}
