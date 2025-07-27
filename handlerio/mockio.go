package handlerio

import (
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jbarzegar/ron-mod-manager/appconfig"
	"github.com/jbarzegar/ron-mod-manager/archive"
)

type MockMod struct {
	Installed bool
	Version   string
	UUID      uuid.UUID
}

// MockIOHandler is a stubbed version of an IOHandler
type MockIOHandler struct {
	Config        *appconfig.AppConfig
	MockedChoices []archive.Choice
	Installed     map[string]MockMod
}

// AddMod will return the mocked choices assigned on struct init
// we assume that `NewMockIoHandler` is used once per test so can assume choices provided are for a specific "mod"
func (h *MockIOHandler) AddMod(archivePath string, outputPath string) ([]archive.Choice, error) {
	defaultChoices := []archive.Choice{{
		Name:     "test-choice",
		FullPath: "/path/to/choice",
	}}
	fmt.Println(slices.Concat(defaultChoices, h.MockedChoices))
	return slices.Concat(defaultChoices, h.MockedChoices), nil
}

func (h *MockIOHandler) InstallMod(payload Installable) error {
	h.Installed[payload.Mod.Name] = MockMod{
		Installed: true,
		Version:   "0.0.0",
		UUID:      *payload.Mod.ActiveVersion,
	}

	return nil
}
