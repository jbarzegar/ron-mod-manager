package handlerio

import (
	"slices"

	"github.com/jbarzegar/ron-mod-manager/archive"
)

// MockIOHandler is a stubbed version of an IOHandler
type MockIOHandler struct {
	mockedChoices []archive.Choice
}

func NewMockIOHandler(archiveChoices []archive.Choice) MockIOHandler {
	return MockIOHandler{mockedChoices: archiveChoices}
}

// AddMod will return the mocked choices assigned on struct init
// we assume that `NewMockIoHandler` is used once per test so can assume choices provided are for a specific "mod"
func (h *MockIOHandler) AddMod(archivePath string, outputPath string) ([]archive.Choice, error) {
	defaultChoices := []archive.Choice{{
		Name:     "test-choice",
		FullPath: "/path/to/choice",
	}}
	return slices.Concat(defaultChoices, h.mockedChoices), nil
}
