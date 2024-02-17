package paths

import (
	"path"

	"github.com/jbarzegar/ron-mod-manager/config"
)

// Get absolute path to games "paks" directory, where .pak bundled mods should be placed
func PaksDir() string {
	conf := config.GetConfig()

	paksDir := path.Join(conf.GameDir, "ReadyOrNot", "Content", "Paks")

	return paksDir
}
