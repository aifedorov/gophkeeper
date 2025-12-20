package binary

import (
	"fmt"

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
	// TODO(human): Implement List method in domain service
	return fmt.Errorf("cli: list command not yet implemented - see TODO(human) in list.go")
}
