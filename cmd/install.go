/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"errors"
	"fmt"
	"strconv"

	statemanagement "github.com/jbarzegar/ron-mod-manager/state-management"
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
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("install called")
		// conf := config.GetConfig()
		// archiveDir := path.Join(conf.ModDir, "archives")
		archives := statemanagement.GetState().Archives

		s, err := strconv.Atoi(args[0])

		if err != nil {
			panic(err)
		}

		selection := s - 1
		// fmt.Println()

		a := utils.FormatArchiveName(archives[selection].FileName)
		fmt.Println("installing " + a)
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
