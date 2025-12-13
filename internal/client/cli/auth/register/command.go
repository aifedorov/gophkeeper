package register

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	"github.com/aifedorov/gophkeeper/internal/client/domain/auth/interfaces"
	"github.com/aifedorov/gophkeeper/pkg/validator"
	"github.com/spf13/cobra"
)

type credentials struct {
	login    string
	password string
}

var creds = &credentials{}

type Command struct {
	cmd     *cobra.Command
	authSrv auth.Service
}

func NewCommand(authSrv auth.Service) (*Command, error) {
	c := &Command{
		authSrv: authSrv,
	}

	cmd := &cobra.Command{
		Use:   "register -l <login> -p <password>",
		Short: "Register a new auth",
		Long:  `Register a new auth with the given login and password.`,
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

func (c *Command) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *Command) run(cmd *cobra.Command, _ []string) error {
	if err := validator.ValidateLogin(creds.login); err != nil {
		return fmt.Errorf("invalid login: %w", err)
	}
	if err := validator.ValidatePassword(creds.password); err != nil {
		return fmt.Errorf("invalid password: %w", err)
	}

	if err := c.authSrv.Register(cmd.Context(), interfaces.Credentials{
		Login:    creds.login,
		Password: creds.password,
	}); err != nil {
		return fmt.Errorf("failed to register: %w", err)
	}

	fmt.Println("You have been successfully logged in")
	return nil
}
