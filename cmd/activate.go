/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/jbarzegar/ron-mod-manager/config"
	statemanagement "github.com/jbarzegar/ron-mod-manager/state-management"
	"github.com/spf13/cobra"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	choices  []string
	cursor   int
	selected map[int]string
}

func initialModel() model {
	state := statemanagement.GetState()

	var choices []string
	for _, mod := range state.Mods {
		choices = append(choices, mod.Name)
	}

	return model{
		choices:  choices,
		selected: make(map[int]string),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	// Is it a key press?
	case tea.KeyMsg:

		// Cool, what was the actual key pressed?
		switch msg.String() {

		// These keys should exit the program.
		case "ctrl+c", "q":
			m.selected = map[int]string{}
			return m, tea.Quit

		// The "up" and "k" keys move the cursor up
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}

		// The "down" and "j" keys move the cursor down
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		// The "enter" key and the spacebar (a literal space) toggle
		// the selected state for the item that the cursor is pointing at.
		case " ":
			_, ok := m.selected[m.cursor]
			if ok {
				delete(m.selected, m.cursor)
			} else {
				m.selected[m.cursor] = m.choices[m.cursor]
			}

		case "enter":
			return m, tea.Quit
		}

	}

	// Return the updated model to the Bubble Tea runtime for processing.
	// Note that we're not returning a command.
	return m, nil
}

func (m model) View() string {
	// The header
	s := "Select mods you want to activate?\n\n"

	// Iterate over our choices
	for i, choice := range m.choices {

		// Is the cursor pointing at this choice?
		cursor := " " // no cursor
		if m.cursor == i {
			cursor = ">" // cursor!
		}

		// Is this choice selected?
		checked := " " // not selected
		if _, ok := m.selected[i]; ok {
			checked = "x" // selected!
		}

		// Render the row
		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
	}

	// The footer
	s += "\nPress q to quit.\n"

	// Send the UI for rendering
	return s
}

// activeCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetConfig()
		p := tea.NewProgram(initialModel())

		// Run select get mods to install
		m, err := p.Run()

		modsToInstall := m.(model).selected

		if err != nil {

			fmt.Printf("Alas, there's been an error: %v", err)
			os.Exit(1)
		}

		// Do install when mods selected
		if len(modsToInstall) >= 1 {
			fmt.Println("activating mods ")

			gamePakDir := path.Join(conf.GameDir, "ReadyOrNot", "Content", "Paks")

			for _, modName := range modsToInstall {
				fmt.Println(modName)

				absModPath, _ := filepath.Abs(conf.ModDir)

				modPath := path.Join(absModPath, "mods", modName)
				m, _ := filepath.Glob(path.Join(modPath, "/*.pak"))

				if len(m) == 1 {
					mod := m[0]
					t := strings.Split(mod, path.Join(modPath)+"/")[1]

					_, err := os.Stat(path.Join(gamePakDir, t))

					state := statemanagement.GetState()
					if !os.IsNotExist(err) {
						// Search all mods and see if the new mod was installed already
						// Account for mods that were installed by ron-mm
						for _, q := range state.Mods {
							if q.Name == modName && q.State == "active" {
								fmt.Println("Mod already installed")
								return
							}
						}
						// Mods may be installed manually and need to be accounted for
						fmt.Println("Mod installed outside of ron-mm's delete mod prior to activating")
					} else {
						for i, q := range state.Mods {
							if q.Name == modName {
								state.Mods[i].State = "active"
							}
						}

						statemanagement.WriteState(state, config.GetConfig())

						fmt.Println(mod)
						fmt.Println(path.Join(gamePakDir, t))
						err = os.Symlink(mod, path.Join(gamePakDir, t))

						if err != nil {
							log.Fatal("why", err)
						}

					}
				}
			}
		} else {
			fmt.Println("No mods to install")
		}

	},
}

func init() {
	rootCmd.AddCommand(activateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// activeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// activeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
