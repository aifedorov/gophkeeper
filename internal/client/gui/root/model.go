package root

import (
	"github.com/aifedorov/gophkeeper/internal/client/container"
	"github.com/aifedorov/gophkeeper/internal/client/gui/auth"
	"github.com/aifedorov/gophkeeper/internal/client/gui/menu"
	tea "github.com/charmbracelet/bubbletea"
)

type screen int

const (
	screenMenu screen = iota
	screenAuth
	screenCategories
)

type Model struct {
	services *container.Services

	currentScreen screen

	menuModel menu.Model
	authModel auth.Model
}

func NewModel(services *container.Services) tea.Model {
	return Model{
		services:      services,
		currentScreen: screenMenu,
		menuModel:     menu.NewModel(),
		authModel:     auth.NewModel(services.AuthSrv),
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
