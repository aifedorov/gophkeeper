package auth

import (
	"context"
	"time"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type NavigateToMenuMsg struct{}

type authSuccessMsg struct {
	session *auth.Session
}

type authErrorMsg struct {
	err error
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.loading && msg.Type != tea.KeyCtrlC {
			break
		}

		switch msg.Type {
		case tea.KeyEnter:
			m.validateField(m.focused)
			if m.focused == len(m.inputs)-1 && m.allFieldsValid() {
				m.loading = true
				m.err = nil
				return m, tea.Batch(
					m.spinner.Tick,
					m.performAuth(),
				)
			}
			m.focused = nextInput(m)
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyTab, tea.KeyCtrlN, tea.KeyDown:
			m.validateField(m.focused)
			m.focused = nextInput(m)
		case tea.KeyShiftTab, tea.KeyCtrlP, tea.KeyUp:
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

	case authSuccessMsg:
		m.loading = false
		m.loggedIn = true
		m.err = nil
		// TODO: navigate to user data screen
		return m, nil

	case authErrorMsg:
		m.loading = false
		m.err = msg.err
		return m, nil

	case spinner.TickMsg:
		if m.loading {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}
	}

	if !m.loading {
		for i := range m.inputs {
			var cmd tea.Cmd
			m.inputs[i], cmd = m.inputs[i].Update(msg)
			cmds = append(cmds, cmd)
		}
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

func (m Model) performAuth() tea.Cmd {
	login := m.inputs[login].Value()
	password := m.inputs[password].Value()
	isNewUser := m.NewUser
	authService := m.authService

	return func() tea.Msg {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		creds := auth.Credentials{
			Login:    login,
			Password: password,
		}

		var session *auth.Session
		var err error

		if isNewUser {
			err = authService.Register(ctx, creds)
		} else {
			err = authService.Login(ctx, creds)
		}

		if err != nil {
			return authErrorMsg{err: err}
		}
		return authSuccessMsg{session: session}
	}
}
