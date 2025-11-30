package menu

import (
	"fmt"
)

func (m Model) View() string {

	// The header
	s := "Welcome to Gophkeeper!\n\n"

	// The menu
	for i, item := range m.Menu {
		cursor := " "
		if m.Cursor == i {
			cursor = ">"
		}

		switch item {
		case Login:
			s += fmt.Sprintf("%s Login\n", cursor)
		case Register:
			s += fmt.Sprintf("%s Register\n", cursor)
		}
	}

	// The footer
	s += "\nctrl+c to quit.\n"

	return s
}
