package app

import (
	"log"

	"github.com/aifedorov/gophkeeper/internal/client/menu"
	"github.com/aifedorov/gophkeeper/internal/client/register"
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
	case screenRegister:
		return m.updateRegister(msg)
	}

	return m, nil
}

func (m Model) updateMenu(msg tea.Msg) (Model, tea.Cmd) {
	log.Printf("menu screen, message type: %T", msg)

	var cmd tea.Cmd
	m.menuModel, cmd = m.menuModel.Update(msg)

	switch msg.(type) {
	case menu.NavigateToRegisterMsg:
		m.currentScreen = screenRegister
		m.registerModel = register.InitialModel()
		return m, cmd
	}

	return m, cmd
}

func (m Model) updateRegister(msg tea.Msg) (Model, tea.Cmd) {
	log.Printf("register screen, message type: %T", msg)

	var cmd tea.Cmd
	m.registerModel, cmd = m.registerModel.Update(msg)

	switch msg.(type) {
	case register.NavigateToMenuMsg:
		m.currentScreen = screenMenu
		return m, cmd
	}

	return m, cmd
}
