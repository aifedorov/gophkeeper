package root

func (m Model) View() string {
	switch m.currentScreen {
	case screenMenu:
		return m.menuModel.View()
	case screenAuth:
		return m.authModel.View()
	}
	return ""
}
