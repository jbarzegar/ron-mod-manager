/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"os"
	"path"
	"strings"

	"github.com/jbarzegar/ron-mod-manager/components"
	"github.com/jbarzegar/ron-mod-manager/config"
	"github.com/jbarzegar/ron-mod-manager/paths"
	s "github.com/jbarzegar/ron-mod-manager/state-management"
	"github.com/spf13/cobra"
)

// deactivateCmd represents the deactivate command
var deactivateCmd = &cobra.Command{
	Use:   "deactivate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		state := s.GetState()
		// Show select to determine which mods to deactivate.
		// Choices will only display active mods in the view
		activeMods := s.GetModsByState("active")

		var choices []string
		for _, m := range activeMods {
			choices = append(choices, m.Name)
		}

		modsToDeactivate := components.SelectMod(choices)
		paksDir := paths.PaksDir()

		for mIdx, m := range modsToDeactivate {
			// Get mod out of state
			mod, err := s.GetModByName(m)

			if err != nil {
				log.Fatal(err)
			}

			// remove symlink (if it's there)
			for _, p := range mod.Paks {
				// s := strings.Split(p, path.Join(config.GetConfig().ModDir, mod.Name, "mods")+"/")[0]

				s := strings.Split(p, path.Join(config.GetConfig().ModDir, "mods", mod.Name)+"/")[1]
				// check if symlink is in dir
				dir := path.Join(paksDir, s)

				_, err := os.Lstat(dir)

				if !os.IsNotExist(err) {
					err := os.Remove(dir)

					if err != nil {
						log.Fatal(err)
					}
				}

			}

			// Update state to signify the mod is inactive
			state.Mods[mIdx].State = "inactive"

			s.WriteState(state, config.GetConfig())
		}
	},
}

func init() {
	rootCmd.AddCommand(deactivateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deactivateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deactivateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
