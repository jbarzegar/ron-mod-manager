package paths

import (
	"path"
	"path/filepath"

	"github.com/jbarzegar/ron-mod-manager/config"
)

func AbsModsDir() string {
	conf := config.GetConfig()
	absModPath, _ := filepath.Abs(conf.ModDir)

	return path.Join(absModPath, "mods")

}

// func ModsDir() string {}
