package root

import (
	"context"
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/cli/auth/login"
	"github.com/aifedorov/gophkeeper/internal/client/cli/auth/register"
	"github.com/aifedorov/gophkeeper/internal/client/cli/commands"
	"github.com/aifedorov/gophkeeper/internal/client/cli/secret/credential"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	domaincredential "github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	clientversion "github.com/aifedorov/gophkeeper/internal/client/version"
	"github.com/spf13/cobra"
)

type Command struct {
	cmd *cobra.Command
}

func NewCommand(authSrv auth.Service, credentialSrv domaincredential.Service) (*Command, error) {
	cmd := &cobra.Command{
		Use:     "gophkeeper",
		Short:   "GophKeeper is a secure password manager",
		Long:    `GophKeeper is a secure password manager that allows you to store and retrieve your passwords securely.`,
		Version: clientversion.Short(),
	}
	cmd.SetVersionTemplate(clientversion.Info() + "\n")

	loginCmd, err := login.NewCommand(authSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create login command: %w", err)
	}
	cmd.AddCommand(loginCmd.GetCommand())

	registerCmd, err := register.NewCommand(authSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create register command: %w", err)
	}
	cmd.AddCommand(registerCmd.GetCommand())

	allCommandsCmd := commands.NewAllCommandsCommand()
	cmd.AddCommand(allCommandsCmd)

	credentialCmd, err := credential.NewCommand(credentialSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create credential command: %w", err)
	}
	cmd.AddCommand(credentialCmd.GetCommand())

	return &Command{
		cmd: cmd,
	}, nil
}

func (r *Command) ExecuteContext(ctx context.Context) error {
	return r.cmd.ExecuteContext(ctx)
}
