package root

import (
	"log"

	auth2 "github.com/aifedorov/gophkeeper/internal/client/gui/auth"
	menu2 "github.com/aifedorov/gophkeeper/internal/client/gui/menu"
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
	m.menuModel = updated.(menu2.Model)

	switch msg.(type) {
	case menu2.NavigateToRegisterMsg:
		m.currentScreen = screenAuth
		m.authModel = auth2.NewModel(m.services.AuthSrv)
		m.authModel.NewUser = true
		return m, cmd
	case menu2.NavigateToLoginMsg:
		m.currentScreen = screenAuth
		m.authModel = auth2.NewModel(m.services.AuthSrv)
		m.authModel.NewUser = false
		return m, cmd
	}

	return m, cmd
}

func (m Model) updateAuth(msg tea.Msg) (Model, tea.Cmd) {
	log.Printf("auth screen, message type: %T", msg)

	var cmd tea.Cmd
	updated, cmd := m.authModel.Update(msg)
	m.authModel = updated.(auth2.Model)

	switch msg.(type) {
	case auth2.NavigateToMenuMsg:
		m.currentScreen = screenMenu
		return m, cmd
	}

	return m, cmd
}
