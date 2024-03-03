package manager

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/jbarzegar/ron-mod-manager/db"
	"github.com/jbarzegar/ron-mod-manager/ent/mod"
	"github.com/jbarzegar/ron-mod-manager/paths"
)

type pakToDisable struct {
	id  int
	mId int
	dir string
}

func Deactivate(modsToDeactivate map[int]string) {
	paksDir := paths.PaksDir()

	for _, modName := range modsToDeactivate {
		// Get mod out of state
		m, err := db.Client().
			Mod.
			Query().
			Where(mod.Name(modName)).Only(context.Background())

		if err != nil {
			log.Fatal(err)
		}

		paks := m.QueryPaks().AllX(context.Background())

		// var dir string
		var paksToDisable []pakToDisable

		for _, p := range paks {
			if p.Installed {
				splitPak := strings.Split(p.Name, string(os.PathSeparator))
				_x := splitPak[len(splitPak)-1]
				paksToDisable = append(paksToDisable, pakToDisable{id: p.ID, dir: path.Join(paksDir, _x), mId: m.ID})
			}
		}

		for _, pd := range paksToDisable {

			_, err = os.Lstat(pd.dir)

			// check if symlink is in dir
			if !os.IsNotExist(err) {
				// remove symlink (if it's there)
				err := os.Remove(pd.dir)

				fmt.Println("Symlink removed")

				if err != nil {
					log.Fatal(err)
				}
			} else {
				fmt.Println("Path doesn't exist: " + pd.dir)
			}

			// Update state to signify the mod is inactive
			db.Client().Mod.Update().Where(mod.ID(pd.mId)).SetState("inactive").Save(context.Background())
		}
	}

}
