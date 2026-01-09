package actions

import (
	"github.com/google/uuid"
	"github.com/jbarzegar/ron-mod-manager/ent"
)

type AddArchiveRequest struct {
	// If ID is supplied in a request
	// It's assumed that you're "reinstalling"
	// a mod
	ID int `json:"id"`
	// ID & ArchivePath cannot exist within the same request body
	ArchivePath string `json:"archivePath"`
	// Not used if ID is provided
	Name string `json:"name"`
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

type ValidationMessage struct {
	// "Warning" | "Fatal" | "Info"
	Kind    string `json:"kind"`
	Message string `json:"message"`
}

type GetArchivesEntry struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Path string `json:"path"`
	// Installed as a mod
	Installed bool `json:"installed"`
	// Installable differs in install where it satisfies the condition where the
	// mod can either be installed or re-installed
	Installable        bool                `json:"installable"`
	PathExists         bool                `json:"pathExists"`
	ValidationMessages []ValidationMessage `json:"validationMessages"`
}

type GetArchiveRequest struct {
	// Scan io for outstanding archives not currently registered by mod manager
	// Default: False
	Untracked bool `json:"untracked"`
}

type GetArchivesResponse struct {
	UntrackedArchives []GetArchivesEntry `json:"untracked"`
	Archives          []GetArchivesEntry `json:"archives"`
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
