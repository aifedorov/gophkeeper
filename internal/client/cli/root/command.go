package root

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/cli/commands"
	"github.com/aifedorov/gophkeeper/internal/client/cli/login"
	"github.com/aifedorov/gophkeeper/internal/client/cli/register"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	clientversion "github.com/aifedorov/gophkeeper/internal/client/version"
	"github.com/spf13/cobra"
)

type RootCommand struct {
	cmd     *cobra.Command
	authSrv auth.Service
}

func NewCommand(authSrv auth.Service) (*RootCommand, error) {
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

	listCmd := commands.NewListCommand()
	cmd.AddCommand(listCmd)

	return &RootCommand{
		cmd:     cmd,
		authSrv: authSrv,
	}, nil
}

func (r *RootCommand) Execute() error {
	return r.cmd.Execute()
}
