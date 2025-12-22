package text

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/text"
	"github.com/spf13/cobra"
)

type ViewCommand struct {
	cmd     *cobra.Command
	textSrv text.Service
}

func NewViewCommand(textSrv text.Service) (*ViewCommand, error) {
	c := &ViewCommand{
		textSrv: textSrv,
	}

	var id string

	cmd := &cobra.Command{
		Use:   "view -i <id>",
		Short: "View a text note",
		Long:  `View a text note by ID. Files larger than 100KB require download command.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, id)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID of the text note to view (required)")
	err := cmd.MarkFlagRequired("id")
	if err != nil {
		return nil, fmt.Errorf("`id` flag as required: %w", err)
	}

	c.cmd = cmd

	return c, nil
}

func (c *ViewCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *ViewCommand) run(cmd *cobra.Command, id string) error {
	if id == "" {
		return fmt.Errorf("cli: id is required")
	}

	content, err := c.textSrv.View(cmd.Context(), id)
	if err != nil {
		return fmt.Errorf("cli: failed to view text note: %w", err)
	}

	fmt.Println(content)

	return nil
}
