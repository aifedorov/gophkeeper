package card

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/card"
	"github.com/spf13/cobra"
)

type DeleteCommand struct {
	cmd     *cobra.Command
	cardSrv card.Service
}

func NewDeleteCommand(cardSrv card.Service) (*DeleteCommand, error) {
	c := &DeleteCommand{
		cardSrv: cardSrv,
	}

	var id string
	cmd := &cobra.Command{
		Use:   "delete -d <id>",
		Short: "Delete a card by ID",
		Long:  `Delete a card by ID`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, id)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "d", "", "ID")
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
	err := c.cardSrv.Delete(cmd.Context(), id)
	if err != nil {
		return fmt.Errorf("cli: failed to delete card: %w", err)
	}

	fmt.Println("Card deleted successfully")
	return nil
}
