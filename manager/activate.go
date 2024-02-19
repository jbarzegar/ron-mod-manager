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

func linkPak(m *types.ModInstall, p types.Pak, modName string) {
	state := statemanagement.GetState()
	absPakModPath := path.Join(paths.AbsModsDir(), m.Name, p.Name)
	absPakGamePath := paths.PaksDir()

	_, err := os.Stat(path.Join(absPakGamePath, p.Name))

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

				for ii, o := range state.Mods[i].Paks {
					if p.Name == o.Name {
						state.Mods[i].Paks[ii].Installed = true
					}

				}
				break
			}
		}

		statemanagement.WriteState(state, config.GetConfig())
		err = os.Symlink(absPakModPath, path.Join(absPakGamePath, p.Name))

		if err != nil {
			log.Fatal("why", err)
		}

	}

}

func Activate(modsToActivate map[int]string) {
	fmt.Println("activating mods ")

	for _, modName := range modsToActivate {
		m, _, _ := statemanagement.GetModByName(modName)

		if len(m.Paks) == 1 {
			p := m.Paks[0]
			linkPak(m, p, modName)

		} else {
			var paks []string
			for _, p := range m.Paks {
				paks = append(paks, p.Name)
			}

			choices := components.SelectMod(paks)
			state := statemanagement.GetState()

			// TODO: linkPak(m, p, modName) should handle recursive mod installs
			for _, p := range choices {
				absPakModPath := path.Join(paths.AbsModsDir(), m.Name, p)
				absPakGamePath := paths.PaksDir()

				splitPak := strings.Split(p, string(os.PathSeparator))
				_x := splitPak[len(splitPak)-1]
				x := path.Join(absPakGamePath, _x)
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

							for ii, o := range state.Mods[i].Paks {
								if p == o.Name {
									state.Mods[i].Paks[ii].Installed = true
									break
								}
							}
							break
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
