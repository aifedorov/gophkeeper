package root

import (
	"github.com/aifedorov/gophkeeper/internal/client/cli/login"
	"github.com/aifedorov/gophkeeper/internal/client/cli/register"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	"github.com/spf13/cobra"
)

type RootCommand struct {
	cmd     *cobra.Command
	authSrv auth.Service
}

func NewCommand(authSrv auth.Service) *RootCommand {
	cmd := &cobra.Command{
		Use:   "gophkeeper",
		Short: "GophKeeper is a secure password manager",
		Long:  `GophKeeper is a secure password manager that allows you to store and retrieve your passwords securely.`,
	}

	loginCmd := login.NewCommand(authSrv)
	cmd.AddCommand(loginCmd.GetCommand())

	registerCmd := register.NewCommand(authSrv)
	cmd.AddCommand(registerCmd.GetCommand())

	return &RootCommand{
		cmd:     cmd,
		authSrv: authSrv,
	}
}

func (r *RootCommand) Execute() error {
	return r.cmd.Execute()
}
