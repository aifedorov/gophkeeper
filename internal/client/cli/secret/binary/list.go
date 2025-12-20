package binary

import (
	"fmt"
	"strings"
	"time"

	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/spf13/cobra"
)

type ListCommand struct {
	cmd       *cobra.Command
	binarySrv binary.Service
}

func NewListCommand(binarySrv binary.Service) (*ListCommand, error) {
	c := &ListCommand{
		binarySrv: binarySrv,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all uploaded files",
		Long:  `List all files that have been uploaded to the server.`,
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
	files, err := c.binarySrv.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("cli: failed to list files: %w", err)
	}
	if len(files) == 0 {
		fmt.Println("No files uploaded yet")
		return nil
	}

	fmt.Printf("%-40s %-35s %-30s %-40s %-30s\n", "ID", "NAME", "SIZE", "UPLOAD DATE", "NOTES")
	fmt.Println(strings.Repeat("-", 160))

	for _, file := range files {
		fmt.Printf("%-40s %-35s %-30d %-40s %-30s\n", file.ID(), file.Name(), file.Size(), file.UploadedAt().Format(time.RFC3339), file.Notes())
	}

	return nil
}
