package credential

import (
	"fmt"

	domaincredential "github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	"github.com/spf13/cobra"
)

type Command struct {
	cmd           *cobra.Command
	credentialSrv domaincredential.Service
}

func NewCommand(credentialSrv domaincredential.Service) (*Command, error) {
	c := &Command{
		credentialSrv: credentialSrv,
	}

	cmd := &cobra.Command{
		Use:   "credential",
		Short: "Manage credential",
		Long:  `Manage credential: create, list, get, update, delete`,
	}

	createCommand, err := NewCreateCommand(credentialSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create create command: %w", err)
	}
	cmd.AddCommand(createCommand.GetCommand())

	listCommand, err := NewListCommand(credentialSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create list command: %w", err)
	}
	cmd.AddCommand(listCommand.GetCommand())

	updateCommand, err := NewUpdateCommand(credentialSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create update command: %w", err)
	}
	cmd.AddCommand(updateCommand.GetCommand())

	deleteCommand, err := NewDeleteCommand(credentialSrv)
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
