package register

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestRegisterInitialState(t *testing.T) {
	m := InitialModel()

	if m.focused != login {
		t.Errorf("expected initial focus on login field, got %d", m.focused)
	}

	if len(m.inputs) != 2 {
		t.Errorf("expected 2 input fields, got %d", len(m.inputs))
	}

	if m.loading {
		t.Error("loading should be false initially")
	}
}

func TestRegisterInit(t *testing.T) {
	m := InitialModel()

	cmd := m.Init()
	if cmd == nil {
		t.Error("expected Init to return a command")
	}
}

func TestRegisterInitWithLoading(t *testing.T) {
	m := InitialModel()
	m.loading = true

	cmd := m.Init()
	if cmd == nil {
		t.Error("expected Init to return spinner tick command when loading")
	}
}

func TestRegisterFocused(t *testing.T) {
	m := InitialModel()

	if m.Focused() != login {
		t.Errorf("expected Focused() to return %d, got %d", login, m.Focused())
	}

	m.focused = password
	if m.Focused() != password {
		t.Errorf("expected Focused() to return %d, got %d", password, m.Focused())
	}
}

func TestRegisterNavigationTabForward(t *testing.T) {
	m := InitialModel()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)

	if m.focused != password {
		t.Errorf("expected focus to move to password field, got %d", m.focused)
	}

	if !m.inputs[password].Focused() {
		t.Error("password input should be focused")
	}

	if m.inputs[login].Focused() {
		t.Error("login input should not be focused")
	}
}

func TestRegisterNavigationShiftTabBackward(t *testing.T) {
	m := InitialModel()
	m.focused = password
	m.inputs[password].Focus()
	m.inputs[login].Blur()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyShiftTab})
	m = updated.(Model)

	if m.focused != login {
		t.Errorf("expected focus to move to login field, got %d", m.focused)
	}
}

func TestRegisterNavigationKeyUp(t *testing.T) {
	m := InitialModel()
	m.focused = password
	m.inputs[password].Focus()
	m.inputs[login].Blur()

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyUp})
	m = updated.(Model)

	if m.focused != login {
		t.Errorf("expected focus to move to login field, got %d", m.focused)
	}

	if !m.inputs[login].Focused() {
		t.Error("login input should be focused")
	}

	if m.inputs[password].Focused() {
		t.Error("password input should not be focused")
	}
}

func TestRegisterNavigationKeyDown(t *testing.T) {
	m := InitialModel()
	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = updated.(Model)

	if m.focused != password {
		t.Errorf("expected focus to move to password field, got %d", m.focused)
	}

	if !m.inputs[password].Focused() {
		t.Error("password input should be focused")
	}

	if m.inputs[login].Focused() {
		t.Error("login input should not be focused")
	}
}

func TestRegisterQuitWithCtrlC(t *testing.T) {
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

func TestRegisterValidationEmptyLogin(t *testing.T) {
	m := InitialModel()
	m.inputs[login].SetValue("")
	m = m.validateField(login)

	if m.inputs[login].Err == nil {
		t.Error("expected validation error for empty login")
	}

	if !strings.Contains(m.inputs[login].Err.Error(), "empty") {
		t.Errorf("expected error about empty login, got %q", m.inputs[login].Err.Error())
	}
}

func TestRegisterValidationEmptyPassword(t *testing.T) {
	m := InitialModel()
	m.inputs[password].SetValue("")
	m = m.validateField(password)

	if m.inputs[password].Err == nil {
		t.Error("expected validation error for empty password")
	}
}

func TestRegisterValidationLongLogin(t *testing.T) {
	longLogin := strings.Repeat("a", LoginMaxLength+1)
	err := loginValidator(longLogin)

	if err == nil {
		t.Error("expected validation error for long login")
	}

	if !strings.Contains(err.Error(), "longer") {
		t.Errorf("expected error about long login, got %q", err.Error())
	}
}

func TestRegisterValidationLongPassword(t *testing.T) {
	longPassword := strings.Repeat("a", PasswordMaxLength+1)
	err := passwordValidator(longPassword)

	if err == nil {
		t.Error("expected validation error for long password")
	}

	if !strings.Contains(err.Error(), "longer") {
		t.Errorf("expected error about long password, got %q", err.Error())
	}
}

func TestRegisterValidationValidFields(t *testing.T) {
	m := InitialModel()
	m.inputs[login].SetValue("testuser")
	m.inputs[password].SetValue("testpass")

	m = m.validateField(login)
	m = m.validateField(password)

	if m.inputs[login].Err != nil {
		t.Errorf("expected no error for valid login, got %v", m.inputs[login].Err)
	}

	if m.inputs[password].Err != nil {
		t.Errorf("expected no error for valid password, got %v", m.inputs[password].Err)
	}
}

func TestRegisterAllFieldsValid(t *testing.T) {
	m := InitialModel()
	m.inputs[login].SetValue("testuser")
	m.inputs[password].SetValue("testpass")

	if !m.allFieldsValid() {
		t.Error("expected all fields to be valid")
	}
}

func TestRegisterAllFieldsInvalid(t *testing.T) {
	m := InitialModel()
	m.inputs[login].SetValue("")
	m.inputs[password].SetValue("")

	if m.allFieldsValid() {
		t.Error("expected fields to be invalid")
	}
}

func TestRegisterBackNavigation(t *testing.T) {
	m := InitialModel()
	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyCtrlB})
	m = updatedModel.(Model)

	if cmd == nil {
		t.Fatal("expected command to be returned")
	}

	msg := cmd()
	if _, ok := msg.(NavigateToMenuMsg); !ok {
		t.Error("expected NavigateToMenuMsg")
	}
}

func TestRegisterValidationOnTab(t *testing.T) {
	m := InitialModel()
	m.inputs[login].SetValue("")

	updated, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = updated.(Model)

	if !m.validated[login] {
		t.Error("login field should be validated after tab")
	}
}

func TestRegisterSubmitInvalidForm(t *testing.T) {
	m := InitialModel()
	m.inputs[login].SetValue("")
	m.inputs[password].SetValue("")
	m.focused = password

	updatedModel, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(Model)

	if m.loading {
		t.Error("should not start loading with invalid form")
	}
}

func TestRegisterSubmitValidForm(t *testing.T) {
	m := InitialModel()
	m.inputs[login].SetValue("testuser")
	m.inputs[password].SetValue("testpass")
	m.focused = password

	updatedModel, cmd := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = updatedModel.(Model)

	if !m.loading {
		t.Error("should start loading with valid form")
	}

	if cmd == nil {
		t.Error("expected spinner tick command")
	}
}
