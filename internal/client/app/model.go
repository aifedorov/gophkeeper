package app

import (
	"github.com/aifedorov/gophkeeper/internal/client/auth"
	"github.com/aifedorov/gophkeeper/internal/client/menu"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenMenu screen = iota
	screenAuth
	screenCategories
)

type Model struct {
	currentScreen screen

	menuModel menu.Model
	authModel auth.Model
}

func InitialModel() tea.Model {
	return Model{
		currentScreen: screenMenu,
		menuModel:     menu.InitialModel(),
		authModel:     auth.InitialModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
