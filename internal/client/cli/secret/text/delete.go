package text

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/text"
	"github.com/spf13/cobra"
)

type DeleteCommand struct {
	cmd     *cobra.Command
	textSrv text.Service
}

func NewDeleteCommand(textSrv text.Service) (*DeleteCommand, error) {
	c := &DeleteCommand{
		textSrv: textSrv,
	}

	var id string

	cmd := &cobra.Command{
		Use:   "delete -i <id>",
		Short: "Delete a text note",
		Long:  `Delete a text note by ID from the server.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, id)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID of the text note to delete (required)")
	err := cmd.MarkFlagRequired("id")
	if err != nil {
		return nil, fmt.Errorf("`id` flag as required: %w", err)
	}

	c.cmd = cmd

	return c, nil
}

func (c *DeleteCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *DeleteCommand) run(cmd *cobra.Command, id string) error {
	if id == "" {
		return fmt.Errorf("cli: id is required")
	}

	if err := c.textSrv.Delete(cmd.Context(), id); err != nil {
		return fmt.Errorf("cli: failed to delete text note: %w", err)
	}

	fmt.Println("Text note deleted successfully")

	return nil
}
