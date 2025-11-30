package register

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	login = iota
	password
)

const (
	LoginMaxLength    = 25
	PasswordMaxLength = 16
)

type Model struct {
	inputs    []textinput.Model
	spinner   spinner.Model
	focused   int
	err       error
	validated []bool
	loading   bool
}

func InitialModel() Model {
	inputs := make([]textinput.Model, 2)

	li := textinput.New()
	li.CharLimit = LoginMaxLength
	li.Prompt = "Login: "
	li.Validate = loginValidator
	li.Focus()
	inputs[login] = li

	pi := textinput.New()
	pi.CharLimit = PasswordMaxLength
	pi.Prompt = "Password: "
	pi.Validate = passwordValidator
	inputs[password] = pi

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return Model{
		inputs:    inputs,
		spinner:   s,
		focused:   login,
		err:       nil,
		validated: make([]bool, 2),
		loading:   false,
	}
}

func (m Model) Init() tea.Cmd {
	if m.loading {
		return m.spinner.Tick
	}
	return textinput.Blink
}

func (m Model) allFieldsValid() bool {
	for _, input := range m.inputs {
		err := input.Validate(input.Value())
		if err != nil {
			return false
		}
	}
	return true
}

func (m Model) Focused() int {
	return m.focused
}

func loginValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("login can't be empty")
	}
	if len(s) > LoginMaxLength {
		return fmt.Errorf("login can't be longer than 25 characters")
	}
	return nil
}

func passwordValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("password can't be empty")
	}
	if len(s) > PasswordMaxLength {
		return fmt.Errorf("password can't be longer than 16 characters")
	}
	return nil
}
