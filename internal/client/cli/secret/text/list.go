package text

import (
	"fmt"
	"strings"
	"time"

	"github.com/aifedorov/gophkeeper/internal/client/domain/text"
	"github.com/spf13/cobra"
)

type ListCommand struct {
	cmd     *cobra.Command
	textSrv text.Service
}

func NewListCommand(textSrv text.Service) (*ListCommand, error) {
	c := &ListCommand{
		textSrv: textSrv,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all text notes",
		Long:  `List all text notes that have been uploaded to the server.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd)
		},
	}

	c.cmd = cmd

	return c, nil
}

func (c *ListCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *ListCommand) run(cmd *cobra.Command) error {
	files, err := c.textSrv.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("cli: failed to list text notes: %w", err)
	}
	if len(files) == 0 {
		fmt.Println("No text notes uploaded yet")
		return nil
	}

	fmt.Printf("%-40s %-35s %-30s %-40s %-30s\n", "ID", "NAME", "SIZE", "UPLOAD DATE", "NOTES")
	fmt.Println(strings.Repeat("-", 160))

	for _, file := range files {
		fmt.Printf("%-40s %-35s %-30d %-40s %-30s\n", file.ID(), file.Name(), file.Size(), file.UploadedAt().Format(time.RFC3339), file.Notes())
	}

	return nil
}
