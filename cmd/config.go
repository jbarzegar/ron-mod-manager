/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/jbarzegar/ron-mod-manager/config"
	"github.com/spf13/cobra"
)

// json | text
var format string

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()

		switch format {
		case "json":
			j, err := json.Marshal(conf)

			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(j))

			break
		case "text":
			fmt.Println("\n")
			fmt.Println("Game Directory: ", conf.GameDir)
			fmt.Println("Modding Instance: ", conf.ModDir)
		default:
			log.Fatal(errors.New(fmt.Sprintf("Unsupported format", format)))
		}

	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	configCmd.Flags().StringVar(
		&format, "format", "text",
		"Configure what format the config will be printed as (json or plain text)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
