package handlerio

import (
	"github.com/jbarzegar/ron-mod-manager/archive"
	"github.com/jbarzegar/ron-mod-manager/ent"
)

type InstallableAssets struct {
	Paks []*ent.Pak
}
type Installable struct {
	Mod         *ent.Mod
	Assets      InstallableAssets
	ArchivePath string
	OutPath     string
}

// IOHandler largely mirrors handler.Handler
type IOHandler interface {
	AddArchive(archivePath string, outputPath string) ([]archive.Choice, error)
	InstallMod(payload Installable) error
	UninstallMod(pakPaths []*ent.Pak) error
	DeleteMod(modPath string) error
	PathExists(path string) bool
	GetUntracked(registered []string) ([]*unTracked, error)
}
