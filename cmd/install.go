/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"log"

	"github.com/jbarzegar/ron-mod-manager/components"
	"github.com/jbarzegar/ron-mod-manager/db"
	"github.com/jbarzegar/ron-mod-manager/ent/archive"
	"github.com/jbarzegar/ron-mod-manager/manager"
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
		archives, err := db.Client().
			Archive.
			Query().
			Where(archive.InstalledEQ(false)).
			All(context.Background())

		if err != nil {
			log.Fatal(err)
		}

		var choices []string
		for _, a := range archives {
			choices = append(choices, a.Name)
		}

		selected := components.SelectMod(choices)

		for _, s := range selected {
			manager.Install(s)
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
