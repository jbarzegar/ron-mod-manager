/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strconv"

	"github.com/iancoleman/strcase"
	"github.com/jbarzegar/ron-mod-manager/config"
	statemanagement "github.com/jbarzegar/ron-mod-manager/state-management"
	"github.com/jbarzegar/ron-mod-manager/types"
	"github.com/jbarzegar/ron-mod-manager/utils"
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
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires 1 arg")
		}

		selection, _ := strconv.Atoi(args[0])

		if selection < 1 {
			return errors.New("selection index must be at least `1`")
		}

		return nil
	},
	// TODO: Allow for multiple mods to be installed?
	Run: func(cmd *cobra.Command, args []string) {
		cf := config.GetConfig()
		state := statemanagement.GetState()

		archives := state.Archives

		s, err := strconv.Atoi(args[0])

		if err != nil {
			panic(err)
		}

		// Get the real selection index
		selectionIdx := s - 1
		selection := archives[selectionIdx]

		name := strcase.ToSnake(utils.FormatArchiveName(selection.FileName))

		fmt.Println("installing " + name)

		modDir := path.Join(cf.ModDir, "mods", name)

		// TODO: When a archive is already extracted, prompt for overwrite
		err = utils.ExtractArchive(selection.FileName, modDir, false)

		if err != nil {
			panic(err)
		}

		// List all paks in file
		// TODO: Make this recursive, rn it will only do a shallow check
		matches, _ := filepath.Glob(path.Join(modDir, "*.pak"))

		mod := types.ModInstall{ArchiveName: selection.FileName, Name: name, State: "inactive"}
		mod.Paks = matches

		// Look through state to see if mod is already installed
		installed := false
		for _, x := range state.Mods {
			if x.Name == name {
				installed = true
			}
		}

		if !installed {
			// append mod and write to state
			state.Mods = append(state.Mods, mod)
			statemanagement.WriteState(state, cf)
		} else {
			fmt.Println("Mod already installed")
		}

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
