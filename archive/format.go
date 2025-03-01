package archive

import (
	"path"
	"strings"

	"github.com/gabriel-vasile/mimetype"
)

func SplitModPath(modDir string, p string) string {
	s := strings.Split(p, path.Join(modDir, "mods")+"/")[1]

	return s
}

// Remove absolute file leaving the file name itself
func SplitArchivePath(modDir string, p string) string {
	s := strings.Split(p, path.Join(modDir, "archives")+"/")[1]

	return s
}

// remove path and file extension
func FormatArchiveName(name string) (string, error) {
	ext, err := mimetype.DetectFile(name)
	if err != nil {
		return "", err
	}

	return strings.Split(name, ext.Extension())[0], nil
}
