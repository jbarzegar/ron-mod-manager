package utils

import (
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jbarzegar/ron-mod-manager/config"
)

// remove path and file extension
func FormatArchiveName(name string) string {
	ext, err := mimetype.DetectFile(name)

	if err != nil {
		panic(err)
	}

	s := strings.Split(name, path.Join(config.GetConfig().ModDir, "archives")+"/")[1]

	return strings.Split(s, ext.Extension())[0]

}
