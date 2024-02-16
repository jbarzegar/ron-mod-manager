package utils

import (
	"path"
	"strings"

	"github.com/jbarzegar/ron-mod-manager/config"
)

func FormatArchiveName(name string) string {
	s := strings.Split(name, path.Join(config.GetConfig().ModDir, "archives")+"/")[1]

	return s

}
