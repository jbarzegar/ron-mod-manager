package manager

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jbarzegar/ron-mod-manager/components"
	"github.com/jbarzegar/ron-mod-manager/db"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/mod"
	"github.com/jbarzegar/ron-mod-manager/ent/pak"
	"github.com/jbarzegar/ron-mod-manager/paths"
)

func linkPak(m *ent.Mod, p *ent.Pak, modName string) {
	mods, err := db.Client().Mod.Query().All(context.Background())
	absPakModPath := path.Join(paths.AbsModsDir(), m.Name, p.Name)
	absPakGamePath := paths.PaksDir()

	_, err = os.Stat(path.Join(absPakGamePath, p.Name))

	if !os.IsNotExist(err) {

		if err != nil {
			log.Fatalf("Error getting mods", err)
		}

		// Search all mods and see if the new mod was installed already
		// Account for mods that were installed by ron-mm
		for _, q := range mods {
			if q.Name == modName && q.State == "active" {
				fmt.Println("Mod already installed")
				return
			}
		}
		// Mods may be installed manually and need to be accounted for
		fmt.Println("Mod installed outside of ron-mm's delete mod prior to activating")

	} else {
		for _, q := range mods {
			if q.Name == modName {
				m, err := db.Client().Mod.Query().Where(mod.ID(q.ID)).Only(context.Background())

				if err != nil {
					log.Fatalf("Error getting mod", err)
				}

				updater := m.Update().SetState("active")
				for _, o := range m.QueryPaks().AllX(context.Background()) {
					if p.Name == o.Name {
						o.Update().
							Where(pak.ID(o.ID)).
							SetInstalled(true).
							SaveX(context.Background())
					}
				}

				updater.Save(context.Background())
				break
			}
		}

		err = os.Symlink(absPakModPath, path.Join(absPakGamePath, p.Name))

		if err != nil {
			log.Fatal("why", err)
		}

	}

}

func Activate(modsToActivate map[int]string) {
	fmt.Println("activating mods ")

	for _, modName := range modsToActivate {
		m, err := db.Client().Mod.Query().
			Where(mod.Name(modName)).
			Only(context.Background())

		if err != nil {
			log.Fatalf("Error getting paks", err)
		}

		paks, err := m.QueryPaks().All(context.Background())

		if err != nil {
			log.Fatalf("Error getting paks", err)
		}

		if len(paks) == 1 {
			p := paks[0]
			linkPak(m, p, modName)

		} else {
			var pakNames []string
			for _, p := range paks {
				pakNames = append(pakNames, p.Name)
			}

			choices := components.SelectMod(pakNames)
			mods := db.Client().Mod.Query().AllX(context.Background())

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
					for _, q := range mods {
						if q.Name == modName && q.State == "active" {
							fmt.Println("Mod already installed")
							return
						}
					}
					// Mods may be installed manually and need to be accounted for
					fmt.Println("Mod installed outside of ron-mm's delete mod prior to activating")

				} else {
					for _, q := range mods {
						if q.Name == modName {
							db.Client().Mod.
								Update().
								Where(mod.ID(q.ID)).
								SetState("active").
								Save(context.Background())

							paks := db.Client().Mod.
								Query().
								Where(mod.ID(q.ID)).
								QueryPaks().
								AllX(context.Background())

							for _, o := range paks {
								if p == o.Name {
									fmt.Println(o.ID, o.Name)
									_, err := db.Client().Pak.
										Update().
										Where(pak.ID(o.ID)).
										SetInstalled(true).
										Save(context.Background())
									if err != nil {
										log.Fatal(err)
									}
									break
								}
							}
							break
						}
					}

					err = os.Symlink(absPakModPath, x)

					if err != nil {
						log.Fatal("why", err)
					}

				}

			}
		}

	}

}
