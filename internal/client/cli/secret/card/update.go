package card

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/card"
	"github.com/spf13/cobra"
)

type UpdateCommand struct {
	cmd     *cobra.Command
	cardSrv card.Service
}

func NewUpdateCommand(cardSrv card.Service) (*UpdateCommand, error) {
	c := &UpdateCommand{
		cardSrv: cardSrv,
	}

	cardInput := inputCard{}

	cmd := &cobra.Command{
		Use:   "update -i <id> -n <name> -u <number> -e <expired_date> -o <card_holder_name> -c <cvv> -i <info>",
		Short: "Update a card by ID",
		Long:  `Update a card by ID with the given name, number, expired date, card holder name, cvv and optional notes.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, cardInput)
		},
	}

	cmd.Flags().StringVarP(&cardInput.id, "id", "d", "", "ID")
	cmd.Flags().StringVarP(&cardInput.name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&cardInput.number, "number", "u", "", "Card number")
	cmd.Flags().StringVarP(&cardInput.expiredDate, "expired_date", "e", "", "Expired date")
	cmd.Flags().StringVarP(&cardInput.cardHolderName, "card_holder_name", "o", "", "Card holder name")
	cmd.Flags().StringVarP(&cardInput.cvv, "cvv", "c", "", "CVV")
	cmd.Flags().StringVarP(&cardInput.notes, "info", "i", "", "Info")

	c.cmd = cmd

	return c, nil
}

func (c *UpdateCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *UpdateCommand) run(cmd *cobra.Command, cardInput inputCard) error {
	if err := cardInput.Validate(); err != nil {
		return fmt.Errorf("cli: failed to validate card: %w", err)
	}

	err := c.cardSrv.Update(cmd.Context(), cardInput.id, card.Card{
		ID:             cardInput.id,
		Name:           cardInput.name,
		Number:         cardInput.number,
		ExpiredDate:    cardInput.expiredDate,
		CardHolderName: cardInput.cardHolderName,
		Cvv:            cardInput.cvv,
		Notes:          cardInput.notes,
	})
	if err != nil {
		return fmt.Errorf("cli: failed to update card: %w", err)
	}

	fmt.Println("Card updated successfully")
	return nil
}
