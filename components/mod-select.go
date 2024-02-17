package components

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	statemanagement "github.com/jbarzegar/ron-mod-manager/state-management"
)

type ModModel struct {
	choices  []string
	cursor   int
	selected map[int]string
}

// Setup model by pulling mods from state
// TODO: Make it possible to grab certain mods (ie, active, inactive etc)
func selectModInitialModel(filter string) ModModel {
	state := statemanagement.GetState()

	var choices []string

	for _, mod := range state.Mods {
		switch filter {
		case "active", "inactive":
			if mod.State == filter {
				choices = append(choices, mod.Name)
			}
		case "":
			choices = append(choices, mod.Name)

		// handle unknown filters
		default:
			fmt.Println("WARN unsupported filter:", filter, " handling as if unfiltered")
			choices = append(choices, mod.Name)

		}
		// if filter == nil {
		// }

	}

	return ModModel{
		choices:  choices,
		selected: make(map[int]string),
	}
}

func (m ModModel) Init() tea.Cmd {
	return nil
}

func (m ModModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m ModModel) View() string {
	// The header
	s := "Select mods you want to deactivate\n\n"

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

func SelectMod(filter string) map[int]string {
	p := tea.NewProgram(selectModInitialModel(filter))

	r, err := p.Run()

	if err != nil {
		panic(err)
	}

	return r.(ModModel).selected

}
