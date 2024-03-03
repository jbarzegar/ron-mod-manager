package manager

import (
	"context"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jbarzegar/ron-mod-manager/db"
	"github.com/jbarzegar/ron-mod-manager/ent"
	"github.com/jbarzegar/ron-mod-manager/ent/archive"
	"github.com/jbarzegar/ron-mod-manager/paths"
	"github.com/jbarzegar/ron-mod-manager/types"
	"github.com/jbarzegar/ron-mod-manager/utils"
)

func Install(n string) {
	absArchivePath := path.Join(paths.AbsArchiveDir(), n)

	arh, err := db.Client().
		Archive.
		Query().
		Where(archive.Name(absArchivePath)).
		Only(context.Background())

	if err != nil {
		log.Fatalf("Err fetching archive %s", err)
	}

	if arh == nil {
		log.Fatal("Could not find archive", n)
	}

	_, err = os.Stat(absArchivePath)
	if arh.Installed && !os.IsNotExist(err) {
		fmt.Println("Mod already installed")
		return
	}

	absModPath := path.Join(paths.AbsModsDir(), arh.Name)

	fmt.Println(absArchivePath, absModPath)

	err = utils.ExtractArchive(absArchivePath, absModPath, false)

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

	mod := types.ModInstall{ArchiveName: n,
		Name: arh.Name, Paks: []types.Pak{}, State: "inactive"}

	for _, m := range matches {
		relPakPath := strings.Split(m, absModPath+"/")[1]
		mod.Paks = append(mod.Paks, types.Pak{Name: relPakPath, Installed: false})
	}

	// Install mod
	m, err := db.Client().Mod.Create().
		SetName(mod.Name).
		SetState("inactive").
		SetArchive(arh).
		Save(context.Background())

	if err != nil {
		log.Fatal(err)
	}
	// -----

	// Install Paks from Mod
	_, err = db.Client().
		Pak.
		MapCreateBulk(mod.Paks, func(c *ent.PakCreate, idx int) {
			p := mod.Paks[idx]
			c.SetName(p.Name).SetInstalled(p.Installed).SetMod(m)
		}).Save(context.Background())

	if err != nil {
		log.Fatal(err)
	}
	// -----

	// Update Archive to be installed
	_, err = arh.
		Update().
		Where(archive.ID(arh.ID)).
		SetInstalled(true).
		Save(context.Background())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Mod installed")

	x, _ := db.Client().Mod.Query().All(context.Background())

	for _, w := range x {
		q, _ := w.QueryPaks().All(context.Background())
		for _, qq := range q {
			fmt.Println(qq)
		}
	}

}
