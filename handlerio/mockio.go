package handlerio

import "github.com/jbarzegar/ron-mod-manager/archive"

// MockIOHandler is a stubbed version of an IOHandler
type MockIOHandler struct{}

func NewMockIOHandler() MockIOHandler {
	return MockIOHandler{}
}

// TODO: Figure out how to pass choices during init
func (h *MockIOHandler) AddMod(archivePath string, outputPath string) ([]archive.Choice, error) {
	return []archive.Choice{}, nil
}
