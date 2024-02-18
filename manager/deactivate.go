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

	for mIdx, m := range modsToDeactivate {
		// Get mod out of state
		mod, err := statemanagement.GetModByName(m)

		if err != nil {
			log.Fatal(err)
		}

		// remove symlink (if it's there)
		for _, p := range mod.Paks {

			if p.Installed {
				// check if symlink is in dir
				dir := path.Join(paksDir, p.Name)
				fmt.Println(dir)
				_, err := os.Lstat(dir)
				if !os.IsNotExist(err) {
					err := os.Remove(dir)

					fmt.Println("Symlink removed")

					if err != nil {
						log.Fatal(err)
					}
				} else {
					x := path.Join(strings.Split(p.Name, string(os.PathSeparator))[1])
					log.Fatalf("Path doesn't exist", x)
				}
			}

		}

		// Update state to signify the mod is inactive
		state.Mods[mIdx].State = "inactive"

		statemanagement.WriteState(state, config.GetConfig())

	}

}
