package root

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/container"
	tea "github.com/charmbracelet/bubbletea"
)

type Root struct {
	services *container.Services
}

func NewRoot(services *container.Services) *Root {
	return &Root{services: services}
}

func (a *Root) Run() error {
	logFile, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		return fmt.Errorf("failed to open log file: %w", err)
	}
	defer func() {
		_ = logFile.Close()
	}()

	p := tea.NewProgram(NewModel(a.services))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}
	return nil
}
