package card

import (
	"fmt"
	"strings"

	"github.com/aifedorov/gophkeeper/internal/client/domain/card"
	"github.com/spf13/cobra"
)

type ListCommand struct {
	cmd     *cobra.Command
	cardSrv card.Service
}

func NewListCommand(cardSrv card.Service) (*ListCommand, error) {
	c := &ListCommand{
		cardSrv: cardSrv,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all cards",
		Long:  `List all cards`,
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
	cards, err := c.cardSrv.List(cmd.Context())
	if err != nil {
		return fmt.Errorf("cli: failed to list cards: %w", err)
	}

	if len(cards) == 0 {
		fmt.Println("No cards found")
		return nil
	}

	fmt.Printf("%-36s %-20s %-20s %-12s %-20s %-5s %s\n", "ID", "NAME", "NUMBER", "EXPIRED", "HOLDER", "CVV", "NOTES")
	fmt.Println(strings.Repeat("-", 150))

	for _, c := range cards {
		fmt.Printf("%-36s %-20s %-20s %-12s %-20s %-5s %s\n", c.ID, c.Name, c.Number, c.ExpiredDate, c.CardHolderName, c.Cvv, c.Notes)
	}

	return nil
}
