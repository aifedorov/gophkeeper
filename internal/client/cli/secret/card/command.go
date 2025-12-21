package card

import (
	"fmt"

	domaincard "github.com/aifedorov/gophkeeper/internal/client/domain/card"
	"github.com/spf13/cobra"
)

type Command struct {
	cmd     *cobra.Command
	cardSrv domaincard.Service
}

func NewCommand(cardSrv domaincard.Service) (*Command, error) {
	c := &Command{
		cardSrv: cardSrv,
	}

	cmd := &cobra.Command{
		Use:   "card",
		Short: "Manage card",
		Long:  `Manage card: create, list, get, update, delete`,
	}

	createCommand, err := NewCreateCommand(cardSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create create command: %w", err)
	}
	cmd.AddCommand(createCommand.GetCommand())

	listCommand, err := NewListCommand(cardSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create list command: %w", err)
	}
	cmd.AddCommand(listCommand.GetCommand())

	updateCommand, err := NewUpdateCommand(cardSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create update command: %w", err)
	}
	cmd.AddCommand(updateCommand.GetCommand())

	deleteCommand, err := NewDeleteCommand(cardSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create delete command: %w", err)
	}
	cmd.AddCommand(deleteCommand.GetCommand())

	c.cmd = cmd

	return c, nil
}

func (c *Command) GetCommand() *cobra.Command {
	return c.cmd
}
