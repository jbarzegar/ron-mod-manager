/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"

	"github.com/jbarzegar/ron-mod-manager/components"
	"github.com/jbarzegar/ron-mod-manager/db"
	"github.com/jbarzegar/ron-mod-manager/ent/mod"
	"github.com/jbarzegar/ron-mod-manager/manager"
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
		// Show select to determine which mods to deactivate.
		// Choices will only display active mods in the view
		activeMods := db.Client().Mod.
			Query().
			Where(mod.State("active")).
			AllX(context.Background())

		var choices []string
		for _, m := range activeMods {
			choices = append(choices, m.Name)
		}

		modsToDeactivate := components.SelectMod(choices)

		manager.Deactivate(modsToDeactivate)
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
