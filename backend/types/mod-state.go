package types

type MMConfig struct {
	GameDir string `json:"gameDir"`
	ModDir  string `json:"modDir"`
}

type Pak struct {
	Name      string `json:"name"`
	Installed bool   `json:"installed"`
}

type ModInstall struct {
	Name        string `json:"name"`
	ArchiveName string `json:"archiveName"`
	// active | inactive
	State string `json:"state"`
	Paks  []Pak  `json:"paks"`
}

type Archive struct {
	Md5Sum      string `json:"md5Sum"`
	Name        string `json:"name"`
	ArchiveFile string `json:"archiveName"`
	Installed   bool   `json:"installed"`
}

type ModState struct {
	Mods       []ModInstall `json:"mods"`
	Archives   []Archive    `json:"archives"`
	ArchiveSum string       `json:"archiveSum"`
}
