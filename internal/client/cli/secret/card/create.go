package card

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/card"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CreateCommand struct {
	cmd     *cobra.Command
	cardSrv card.Service
}

func NewCreateCommand(cardSrv card.Service) (*CreateCommand, error) {
	c := &CreateCommand{
		cardSrv: cardSrv,
	}

	cardInput := inputCard{}

	cmd := &cobra.Command{
		Use:   "create -n <name> -u <number> -e <expired_date> -o <card_holder_name> -c <cvv> [-i <info>]",
		Short: "Create a new card",
		Long:  `Create a new card with the given name, number, expired date, card holder name, cvv and optional notes.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, cardInput)
		},
	}

	cmd.Flags().StringVarP(&cardInput.name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&cardInput.number, "number", "u", "", "Card number")
	cmd.Flags().StringVarP(&cardInput.expiredDate, "expired_date", "e", "", "Expired date")
	cmd.Flags().StringVarP(&cardInput.cardHolderName, "card_holder_name", "o", "", "Card holder name")
	cmd.Flags().StringVarP(&cardInput.cvv, "cvv", "c", "", "CVV")
	cmd.Flags().StringVarP(&cardInput.notes, "info", "i", "", "Info")

	c.cmd = cmd

	return c, nil
}

func (c *CreateCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *CreateCommand) run(cmd *cobra.Command, cardInput inputCard) error {
	id := uuid.New().String()
	newCard, err := card.NewCard(id, cardInput.name, cardInput.number, cardInput.expiredDate, cardInput.cardHolderName, cardInput.cvv, cardInput.notes)
	if err != nil || newCard == nil {
		return fmt.Errorf("cli: failed to create card: %w", err)
	}

	if err := newCard.Validate(); err != nil {
		return fmt.Errorf("cli: failed to validate card: %w", err)
	}

	if err := c.cardSrv.Create(cmd.Context(), *newCard); err != nil {
		return fmt.Errorf("cli: failed to create card: %w", err)
	}

	fmt.Println("Card created successfully")
	return nil
}
