package handlerio

import (
	"github.com/jbarzegar/ron-mod-manager/archive"
	"github.com/jbarzegar/ron-mod-manager/ent"
)

type InstallableAssets struct {
	Pak []*ent.Pak
}
type Installable struct {
	Mod         *ent.Mod
	Assets      InstallableAssets
	ArchivePath string
	OutPath     string
}

// IOHandler largely mirrors handler.Handler
type IOHandler interface {
	AddMod(archivePath string, outputPath string) ([]archive.Choice, error)
	InstallMod(payload Installable) error
}
