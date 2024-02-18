package paths

import (
	"path"
	"path/filepath"

	"github.com/jbarzegar/ron-mod-manager/config"
)

func absBasePath() string {
	conf := config.GetConfig()
	absModPath, _ := filepath.Abs(conf.ModDir)

	return absModPath
}

func AbsArchiveDir() string {
	return path.Join(absBasePath(), "archives")
}

func AbsModsDir() string {
	return path.Join(absBasePath(), "mods")
}
