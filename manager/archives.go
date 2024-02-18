package manager

// func GetArchiveByName(name string) (*types.Archive, error) {
// 	state := statemanagement.GetState()

// 	var selectedMod *types.ModInstall = nil
// 	for _, m := range state.Mods {
// 		if m.Name == name {
// 			selectedMod = &m
// 			break
// 		}
// 	}

// 	if selectedMod == nil {
// 		return nil, errors.New("Mod not found")
// 	}

// 	return selectedMod, nil
// }

// func MakeArchive(name string) types.Archive {
// 	state := statemanagement.GetState()

// 	mod, _ := statemanagement.GetModByName(name)

// 	mod.
// 		fmt.Println("Wut", archive.FileName, archive.Installed)

// 	return archive
// }

// func ListArchives() string[] {

// }
