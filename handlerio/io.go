package handlerio

import (
	"github.com/jbarzegar/ron-mod-manager/archive"
)

// IOHandler largely mirrors handler.Handler
type IOHandler interface {
	AddMod(archivePath string, outputPath string) ([]archive.Choice, error)
}
