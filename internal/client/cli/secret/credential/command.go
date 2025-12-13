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
		Short: "Manage credentials",
		Long:  `Manage credentials: create, list, get, update, delete`,
	}

	createCommand, err := NewCreateCommand(credentialSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create create command: %w", err)
	}

	cmd.AddCommand(createCommand.GetCommand())

	c.cmd = cmd

	return c, nil
}

func (c *Command) GetCommand() *cobra.Command {
	return c.cmd
}
