package menu

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestMenuInitialState(t *testing.T) {
	m := InitialModel()

	if m.Cursor != 0 {
		t.Errorf("expected initial cursor to be 0, got %d", m.Cursor)
	}

	if len(m.Menu) != 2 {
		t.Errorf("expected menu to have 2 items, got %d", len(m.Menu))
	}
}

func TestMenuInit(t *testing.T) {
	m := InitialModel()

	cmd := m.Init()
	if cmd != nil {
		t.Error("expected Init to return nil for menu")
	}
}

func TestMenuBoundaryDown(t *testing.T) {
	m := InitialModel()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)
	updated, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)

	if m.Cursor != 1 {
		t.Errorf("cursor should stay at last position, got %d", m.Cursor)
	}
}

func TestMenuBoundaryUp(t *testing.T) {
	m := InitialModel()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)

	if m.Cursor != 0 {
		t.Errorf("cursor should stay at first position, got %d", m.Cursor)
	}
}

func TestMenuSelectLogin(t *testing.T) {
	m := InitialModel()

	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(Model)

	if cmd == nil {
		t.Fatal("expected command to be returned")
	}

	msg := cmd()
	if _, ok := msg.(NavigateToLoginMsg); !ok {
		t.Error("expected NavigateToLoginMsg")
	}
}

func TestMenuSelectRegister(t *testing.T) {
	m := InitialModel()
	m.Cursor = 1
	m.Selected = Register

	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(Model)

	if cmd == nil {
		t.Fatal("expected command to be returned")
	}

	msg := cmd()
	if _, ok := msg.(NavigateToRegisterMsg); !ok {
		t.Error("expected NavigateToRegisterMsg")
	}
}

func TestMenuNavigationUp(t *testing.T) {
	m := InitialModel()
	m.Cursor = 1
	m.Selected = Register

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updatedModel.(Model)

	if m.Cursor != 0 {
		t.Errorf("expected cursor to be 0, got %d", m.Cursor)
	}

	if m.Selected != Login {
		t.Errorf("expected selected to be Login, got %v", m.Selected)
	}
}

func TestMenuNavigationDown(t *testing.T) {
	m := InitialModel()

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updatedModel.(Model)

	if m.Cursor != 1 {
		t.Errorf("expected cursor to be 1, got %d", m.Cursor)
	}

	if m.Selected != Register {
		t.Errorf("expected selected to be Register, got %v", m.Selected)
	}
}

func TestMenuQuitWithCtrlC(t *testing.T) {
	m := InitialModel()
	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlC})
	m = updatedModel.(Model)

	if cmd == nil {
		t.Fatal("expected quit command")
	}

	msg := cmd()
	if _, ok := msg.(tea.QuitMsg); !ok {
		t.Error("expected QuitMsg")
	}
}
