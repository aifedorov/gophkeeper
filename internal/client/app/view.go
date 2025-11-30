package app

func (m Model) View() string {
	switch m.currentScreen {
	case screenMenu:
		return m.menuModel.View()
	case screenRegister:
		return m.registerModel.View()
	}
	return ""
}
