package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

type NavigateToRegisterMsg struct{}
type NavigateToLoginMsg struct{}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyUp:
			if m.Cursor > 0 {
				m.Cursor--
				m.Selected = m.Menu[m.Cursor]
			}
		case tea.KeyDown:
			if m.Cursor < len(m.Menu)-1 {
				m.Cursor++
				m.Selected = m.Menu[m.Cursor]
			}
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			switch m.Selected {
			case Login:
				return m, func() tea.Msg {
					return NavigateToLoginMsg{}
				}
			case Register:
				return m, func() tea.Msg {
					return NavigateToRegisterMsg{}
				}
			}
		}
	}
	return m, nil
}
