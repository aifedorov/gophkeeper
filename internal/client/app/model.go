package app

import (
	"github.com/aifedorov/gophkeeper/internal/client/menu"
	"github.com/aifedorov/gophkeeper/internal/client/register"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenMenu screen = iota
	screenRegister
	screenLogin
	screenCategories
)

type Model struct {
	currentScreen screen

	menuModel     menu.Model
	registerModel register.Model
}

func InitialModel() tea.Model {
	return Model{
		currentScreen: screenMenu,
		menuModel:     menu.InitialModel(),
		registerModel: register.InitialModel(),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
