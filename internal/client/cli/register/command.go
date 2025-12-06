package register

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/auth"
	"github.com/spf13/cobra"
)

const (
	LoginMinLength    = 3
	LoginMaxLength    = 25
	PasswordMinLength = 3
	PasswordMaxLength = 16
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

func NewCommand(authSrv auth.Service) *RegisterCommand {
	cmd := &cobra.Command{
		Use:   "register -l <login> -p <password>",
		Short: "Register a new user",
		Long:  `Register a new user with the given login and password.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := loginValidator(creds.login); err != nil {
				return fmt.Errorf("invalid login: %w", err)
			}
			if err := passwordValidator(creds.password); err != nil {
				return fmt.Errorf("invalid password: %w", err)
			}
			fmt.Printf("registering user with login: %s and password: %s\n", creds.login, creds.password)
			return nil
		},
	}

	cmd.Flags().StringVarP(&creds.login, "login", "l", "", "Login")
	cmd.Flags().StringVarP(&creds.password, "password", "p", "", "Password")

	cmd.MarkFlagRequired("login")
	cmd.MarkFlagRequired("password")

	return &RegisterCommand{
		cmd:     cmd,
		authSrv: authSrv,
	}
}

func (c *RegisterCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func loginValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("login can't be empty")
	}
	if len(s) < LoginMinLength {
		return fmt.Errorf("login must be at least 3 characters")
	}
	if len(s) > LoginMaxLength {
		return fmt.Errorf("login can't be longer than 25 characters")
	}
	return nil
}

func passwordValidator(s string) error {
	if len(s) == 0 {
		return fmt.Errorf("password can't be empty")
	}
	if len(s) < PasswordMinLength {
		return fmt.Errorf("password must be at least 3 characters")
	}
	if len(s) > PasswordMaxLength {
		return fmt.Errorf("password can't be longer than 16 characters")
	}
	return nil
}
