package auth

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

const (
	red = lipgloss.Color("#FF0025")
)

var (
	errorMessageStyle = lipgloss.NewStyle().Foreground(red)
)

func (m Model) View() string {
	return fmt.Sprint("" +
		m.header() +
		m.inputs[login].View() + "\n" +
		m.validationError(login) +
		m.inputs[password].View() + "\n" +
		m.validationError(password) +
		m.errorState() +
		m.loadingState() +
		m.successState() +
		"\n(ctrl+b to return back)\n",
	)
}

func (m Model) header() string {
	if m.NewUser {
		return "Create a new login and password\n\n"
	}
	return "Enter login and password\n\n"
}

func (m Model) validationError(fieldIdx int) string {
	if m.validated[fieldIdx] {
		if err := m.inputs[fieldIdx].Err; err != nil {
			return errorMessageStyle.Render(err.Error()) + "\n"
		}
	}
	return ""
}

func (m Model) errorState() string {
	if m.err != nil {
		return errorMessageStyle.Render(m.err.Error()) + "\n"
	}
	return ""
}

func (m Model) loadingState() string {
	action := "Logging in..."
	if m.NewUser {
		action = "Registering..."
	}
	if m.loading {
		return fmt.Sprintf("\n%s %s...\n", m.spinner.View(), action)
	}
	return ""
}

func (m Model) successState() string {
	if !m.loggedIn {
		return ""
	}
	return "\nSuccessfully logged in!\n"
}
