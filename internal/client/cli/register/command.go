package register

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	"github.com/spf13/cobra"
)

type credentials struct {
	login    string
	password string
}

var creds = &credentials{}

type RegisterCommand struct {
	cmd     *cobra.Command
	authSrv auth.Service
}

func NewCommand(authSrv auth.Service) (*RegisterCommand, error) {
	c := &RegisterCommand{
		authSrv: authSrv,
	}

	cmd := &cobra.Command{
		Use:   "register -l <login> -p <password>",
		Short: "Register a new user",
		Long:  `Register a new user with the given login and password.`,
		RunE:  c.run,
	}

	cmd.Flags().StringVarP(&creds.login, "login", "l", "", "Login")
	cmd.Flags().StringVarP(&creds.password, "password", "p", "", "Password")

	if err := cmd.MarkFlagRequired("login"); err != nil {
		return nil, fmt.Errorf("failed to mark login flag as required: %w", err)
	}
	if err := cmd.MarkFlagRequired("password"); err != nil {
		return nil, fmt.Errorf("failed to mark password flag as required: %w", err)
	}

	c.cmd = cmd

	return c, nil
}

func (c *RegisterCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *RegisterCommand) run(cmd *cobra.Command, args []string) error {
	if err := auth.ValidateLogin(creds.login); err != nil {
		return fmt.Errorf("invalid login: %w", err)
	}
	if err := auth.ValidatePassword(creds.password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	if err := c.authSrv.Register(cmd.Context(), auth.Credentials{
		Login:    creds.login,
		Password: creds.password,
	}); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	fmt.Println("You have been successfully logged in")
	return nil
}
