package menu

import (
	tea "github.com/charmbracelet/bubbletea"
)

type Item int

const (
	Login Item = iota
	Register
)

type Model struct {
	Menu     []Item
	Cursor   int
	Selected Item
}

func InitialModel() Model {
	return Model{
		Menu: []Item{Login, Register},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}
