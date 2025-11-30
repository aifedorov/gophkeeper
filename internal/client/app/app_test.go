package app

import (
	"testing"

	"github.com/aifedorov/gophkeeper/internal/client/auth"
	"github.com/aifedorov/gophkeeper/internal/client/menu"
	tea "github.com/charmbracelet/bubbletea"
)

func TestAppInitialState(t *testing.T) {
	m := InitialModel().(Model)

	if m.currentScreen != screenMenu {
		t.Errorf("expected initial screen to be menu, got %v", m.currentScreen)
	}
}

func TestAppInit(t *testing.T) {
	m := InitialModel()

	cmd := m.Init()
	if cmd != nil {
		t.Error("expected Init to return nil for app")
	}
}

func TestAppQuitFromAnyScreen(t *testing.T) {
	m := InitialModel()
	_, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})

	if cmd == nil {
		t.Fatal("expected quit command")
	}

	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Error("expected QuitMsg")
	}
}

func TestAppMenuScreenDelegation(t *testing.T) {
	m := InitialModel().(Model)

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updatedModel.(Model)

	if m.menuModel.Cursor != 1 {
		t.Error("menu model should handle key events when on menu screen")
	}
}

func TestAppRegisterScreenDelegation(t *testing.T) {
	m := InitialModel().(Model)
	m.currentScreen = screenAuth
	m.authModel = auth.InitialModel()

	initialFocus := m.authModel.Focused()

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updatedModel.(Model)

	if m.authModel.Focused() == initialFocus {
		t.Error("auth model should handle tab key and change focus")
	}
}

func TestAppScreenTransitions(t *testing.T) {
	tests := []struct {
		name        string
		startScreen screen
		msg         tea.Msg
		wantScreen  screen
	}{
		{
			name:        "menu to auth",
			startScreen: screenMenu,
			msg:         menu.NavigateToRegisterMsg{},
			wantScreen:  screenAuth,
		},
		{
			name:        "auth to menu",
			startScreen: screenAuth,
			msg:         auth.NavigateToMenuMsg{},
			wantScreen:  screenMenu,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := InitialModel().(Model)
			m.currentScreen = tt.startScreen

			if tt.startScreen == screenAuth {
				m.authModel = auth.InitialModel()
			}

			updatedModel, _ := m.Update(tt.msg)
			m = updatedModel.(Model)

			if m.currentScreen != tt.wantScreen {
				t.Errorf("expected screen %v, got %v", tt.wantScreen, m.currentScreen)
			}
		})
	}
}
