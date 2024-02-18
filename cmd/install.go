/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"github.com/jbarzegar/ron-mod-manager/components"
	"github.com/jbarzegar/ron-mod-manager/manager"
	statemanagement "github.com/jbarzegar/ron-mod-manager/state-management"
	"github.com/spf13/cobra"
)

// func isSupportedMimeType (file fs.File) {}

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// TODO: Allow for multiple mods to be installed?
	Run: func(cmd *cobra.Command, args []string) {
		state := statemanagement.GetState()

		var choices []string
		for _, a := range state.Archives {
			choices = append(choices, a.ArchiveFile)
		}

		selected := components.SelectMod(choices)

		for _, s := range selected {
			manager.Install(s)
		}

		// fmt.Println(selected)
		// cf := config.GetConfig()
		// state := statemanagement.GetState()

		// fmt.Println(state.Archives)

		// s, err := strconv.Atoi(args[0])
		// if err != nil {
		// 	panic(err)
		// }

		// // Get the real selection index
		// selection := state.Archives[s-1]

		// name := strcase.ToSnake(utils.FormatArchiveName(selection.FileName))

		// fmt.Println("installing " + name)
		// modDir := path.Join(cf.ModDir, "mods", name)

		// // TODO: When a archive is already extracted, prompt for overwrite
		// err = utils.ExtractArchive(selection.FileName, modDir, false)

		// if err != nil {
		// 	panic(err)
		// }

		// // List all paks in file
		// // TODO: Make this recursive, rn it will only do a shallow check
		// matches, _ := filepath.Glob(path.Join(modDir, "*.pak"))

		// mod := types.ModInstall{ArchiveName: selection.FileName, Name: name, State: "inactive"}
		// mod.Paks = matches

		// // Look through state to see if mod is already installed
		// installed := false
		// for _, x := range state.Mods {
		// 	if x.Name == name {
		// 		installed = true
		// 	}
		// }

		// if !installed {
		// 	// append mod and write to state
		// 	state.Mods = append(state.Mods, mod)
		// 	statemanagement.WriteState(state, cf)
		// } else {
		// 	fmt.Println("Mod already installed")
		// }

	},
}

func init() {
	rootCmd.AddCommand(installCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// installCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// installCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
