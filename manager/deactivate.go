package manager

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jbarzegar/ron-mod-manager/config"
	"github.com/jbarzegar/ron-mod-manager/paths"
	statemanagement "github.com/jbarzegar/ron-mod-manager/state-management"
)

func Deactivate(modsToDeactivate map[int]string) {
	state := statemanagement.GetState()
	paksDir := paths.PaksDir()

	for _, m := range modsToDeactivate {
		// Get mod out of state
		mod, mIdx, err := statemanagement.GetModByName(m)

		if err != nil {
			log.Fatal(err)
		}

		var dir string
		var pIdx int
		if len(mod.Paks) > 1 {
			for pi, p := range mod.Paks {
				if p.Installed {
					splitPak := strings.Split(p.Name, string(os.PathSeparator))
					_x := splitPak[len(splitPak)-1]
					dir = path.Join(paksDir, _x)
					pIdx = pi
				}
			}
		} else {
			p := mod.Paks[0]
			if p.Installed {
				dir = path.Join(paksDir, p.Name)
				pIdx = 0
			}

		}

		_, err = os.Lstat(dir)

		// check if symlink is in dir
		if !os.IsNotExist(err) {
			// remove symlink (if it's there)
			err := os.Remove(dir)

			fmt.Println("Symlink removed")

			if err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Path doesn't exist: " + dir)
		}

		// Update state to signify the mod is inactive
		state.Mods[mIdx].State = "inactive"
		state.Mods[mIdx].Paks[pIdx].Installed = false

		statemanagement.WriteState(state, config.GetConfig())

	}

}
