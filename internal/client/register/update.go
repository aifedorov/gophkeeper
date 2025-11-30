package register

import (
	tea "github.com/charmbracelet/bubbletea"
)

type NavigateToMenuMsg struct{}

type errMsg error

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	cmds := make([]tea.Cmd, len(m.inputs)+1) // spinner

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			m.validateField(m.focused)
			if m.focused == len(m.inputs)-1 && m.allFieldsValid() {
				m.loading = true
				// TODO: Add submission logic here
				return m, m.spinner.Tick
			}
			m.focused = nextInput(m)
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab, tea.KeyCtrlN:
			m.validateField(m.focused)
			m.focused = nextInput(m)
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.validateField(m.focused)
			m.focused = prevInput(m)
		case tea.KeyCtrlB:
			return m, func() tea.Msg {
				return NavigateToMenuMsg{}
			}
		default:
			m.err = nil
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()

	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	if m.loading {
		m.spinner, cmds[len(cmds)-1] = m.spinner.Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m Model) validateField(fieldIdx int) Model {
	value := m.inputs[fieldIdx].Value()
	m.validated[fieldIdx] = true

	switch fieldIdx {
	case login:
		if err := loginValidator(value); err != nil {
			m.inputs[fieldIdx].Err = err
		} else {
			m.inputs[fieldIdx].Err = nil
		}
	case password:
		if err := passwordValidator(value); err != nil {
			m.inputs[fieldIdx].Err = err
		} else {
			m.inputs[fieldIdx].Err = nil
		}
	}
	return m
}

func nextInput(m Model) int {
	return (m.focused + 1) % len(m.inputs)
}

func prevInput(m Model) int {
	prev := m.focused - 1
	if prev < 0 {
		return len(m.inputs) - 1
	}
	return prev
}
