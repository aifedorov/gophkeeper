package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type App struct{}

func NewApp() *App {
	return &App{}
}

func (a *App) Run() error {
	logFile, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		fmt.Println("Error setting up log:", err)
		os.Exit(1)
	}
	defer func() {
		_ = logFile.Close()
	}()

	p := tea.NewProgram(InitialModel())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to run program: %w", err)
	}
	return nil
}
