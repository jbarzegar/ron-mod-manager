package manager

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/jbarzegar/ron-mod-manager/config"
	"github.com/jbarzegar/ron-mod-manager/paths"
	statemanagement "github.com/jbarzegar/ron-mod-manager/state-management"
)

func Activate(modsToActivate map[int]string) {
	state := statemanagement.GetState()
	fmt.Println("activating mods ")

	for _, modName := range modsToActivate {
		m, _ := statemanagement.GetModByName(modName)

		fmt.Println(m.Name, m.Paks)

		if len(m.Paks) == 1 {
			p := m.Paks[0]
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

		} else {
			log.Fatal("Multiple paks found! unsupported for now")

			for _, p := range m.Paks {
				fmt.Println(p)
			}
		}

	}
}
