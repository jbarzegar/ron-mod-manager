package utils

import (
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
	"github.com/jbarzegar/ron-mod-manager/config"
)

func SplitModPath(p string) string {
	s := strings.Split(p, path.Join(config.GetConfig().ModDir, "mods")+"/")[1]

	return s
}

func SplitArchivePath(p string) string {
	s := strings.Split(p, path.Join(config.GetConfig().ModDir, "archives")+"/")[1]

	return s
}

// remove path and file extension
func FormatArchiveName(name string) string {
	ext, err := mimetype.DetectFile(name)

	if err != nil {
		panic(err)
	}

	return strings.Split(name, ext.Extension())[0]

}
