package app

import (
	"log"

	"github.com/aifedorov/gophkeeper/internal/client/auth"
	"github.com/aifedorov/gophkeeper/internal/client/menu"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	log.Printf("current screen: %v, message type: %T", m.currentScreen, msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	switch m.currentScreen {
	case screenMenu:
		return m.updateMenu(msg)
	case screenAuth:
		return m.updateAuth(msg)
	}

	return m, nil
}

func (m Model) updateMenu(msg tea.Msg) (Model, tea.Cmd) {
	log.Printf("menu screen, message type: %T", msg)

	var cmd tea.Cmd
	updated, cmd := m.menuModel.Update(msg)
	m.menuModel = updated.(menu.Model)

	switch msg.(type) {
	case menu.NavigateToRegisterMsg:
		m.currentScreen = screenAuth
		m.authModel = auth.InitialModel()
		m.authModel.NewUser = true
		return m, cmd
	case menu.NavigateToLoginMsg:
		m.currentScreen = screenAuth
		m.authModel = auth.InitialModel()
		m.authModel.NewUser = false
		return m, cmd
	}

	return m, cmd
}

func (m Model) updateAuth(msg tea.Msg) (Model, tea.Cmd) {
	log.Printf("auth screen, message type: %T", msg)

	var cmd tea.Cmd
	updated, cmd := m.authModel.Update(msg)
	m.authModel = updated.(auth.Model)

	switch msg.(type) {
	case auth.NavigateToMenuMsg:
		m.currentScreen = screenMenu
		return m, cmd
	}

	return m, cmd
}
