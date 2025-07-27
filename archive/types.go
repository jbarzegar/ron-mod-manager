package archive

import "github.com/google/uuid"

type Choice struct {
	UUID     uuid.UUID `json:"uuid"`
	Name     string    `json:"name"`
	FullPath string    `json:"fullPath"`
}
