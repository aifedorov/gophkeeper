package binary

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/spf13/cobra"
)

type DeleteCommand struct {
	cmd       *cobra.Command
	binarySrv binary.Service
}

func NewDeleteCommand(binarySrv binary.Service) (*DeleteCommand, error) {
	c := &DeleteCommand{
		binarySrv: binarySrv,
	}

	var id string
	cmd := &cobra.Command{
		Use:   "delete -i <id>",
		Short: "Delete a file by ID",
		Long:  `Delete a file by ID`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, id)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID")
	if err := cmd.MarkFlagRequired("id"); err != nil {
		return nil, fmt.Errorf("cli: failed to mark id flag as required: %w", err)
	}

	c.cmd = cmd

	return c, nil
}

func (c *DeleteCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *DeleteCommand) run(cmd *cobra.Command, id string) error {
	err := c.binarySrv.Delete(cmd.Context(), id)
	if err != nil {
		return fmt.Errorf("cli: failed to delete file: %w", err)
	}

	fmt.Println("File deleted successfully")
	return nil
}
