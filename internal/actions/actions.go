package actions

import (
	"github.com/google/uuid"
	"github.com/jbarzegar/ron-mod-manager/ent"
)

type AddRequest struct {
	ArchivePath string `json:"archivePath"`
	Name        string `json:"name"`
}

type InstallRequest struct {
	ModID      int         `json:"modId"`
	ModVersion string      `json:"modVersionId"`
	Choices    []uuid.UUID `json:"choices"`
}

type StagedMod struct {
	Name              string     `json:"name"`
	State             string     `json:"state"`
	Origin            *string    `json:"origin"`
	ActiveVersion     string     `json:"activeVersion"`
	ActiveVersionUUID *uuid.UUID `json:"activeVersionUUID"`
	Paks              []*ent.Pak `json:"paks"`
}

type StagedResponse struct {
	Mods []StagedMod `json:"mods"`
}

type AllModsEntry struct {
	Name          string            `json:"name"`
	State         string            `json:"state"`
	Origin        *string           `json:"origin"`
	ActiveVersion string            `json:"activeVersion"`
	Versions      []*ent.ModVersion `json:"versions"`
}

type AllModsResponse struct {
	Mods []AllModsEntry `json:"mods"`
}

type UninstallModRequest struct {
	ModIds []int `json:"modIds"`
}

type DeleteModEntry struct {
	ID int `json:"id"`
	// all or active
	DeleteVersions string `json:"deleteVersions"`
	// default: false
	DeleteArchive bool `json:"deleteArchive"`
}

type DeleteModRequest struct {
	Mods []DeleteModEntry `json:"mods"`
}
