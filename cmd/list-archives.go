/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"context"
	"fmt"
	"log"

	"github.com/jbarzegar/ron-mod-manager/db"
	"github.com/spf13/cobra"
)

func handleListArchives(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	archives, err := db.Client().Archive.Query().All(ctx)

	if err != nil {
		log.Fatalf("Err fetching archives", err)
	}

	if len(archives) == 0 {
		fmt.Println("No archives")
	}

	for i, a := range archives {
		var str = a.Name

		if a.Installed {
			str += str + " [Installed]"
		}

		fmt.Println(i+1, str)
	}
}

// listArchivesCmd represents the listArchives command
var listArchivesCmd = &cobra.Command{
	Use:   "list-archives",
	Short: "list mod archives",
	Long:  `Lists mod archives dir`,
	Run:   handleListArchives,
}

func init() {
	rootCmd.AddCommand(listArchivesCmd)
	// rootCmd.Flags().AddFlag()

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listArchivesCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listArchivesCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	// state := stateManagement.GetState()
}
