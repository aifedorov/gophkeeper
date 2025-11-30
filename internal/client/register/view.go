package register

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
		"Create a new login and password\n" +
		m.inputs[login].View() + "\n" +
		m.validationError(login) +
		m.inputs[password].View() + "\n" +
		m.validationError(password) +
		m.errorState() +
		m.loadingState() + "\n" +
		"(ctrl+b to return back)\n",
	)
}

func (m Model) validationError(fieldIdx int) string {
	if m.validated[fieldIdx] {
		if err := m.inputs[fieldIdx].Err; err != nil {
			return errorMessageStyle.Render(err.Error())
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
	if m.loading {
		return fmt.Sprintf("%s Loading...", m.spinner.View())
	}
	return ""
}
