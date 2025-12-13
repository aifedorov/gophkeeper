package credential

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type credentials struct {
	name     string
	login    string
	password string
	notes    string
}

var creds = &credentials{}

type CreateCommand struct {
	cmd           *cobra.Command
	credentialSrv credential.Service
}

func NewCreateCommand(credentialSrv credential.Service) (*CreateCommand, error) {
	c := &CreateCommand{
		credentialSrv: credentialSrv,
	}

	cmd := &cobra.Command{
		Use:   "create -n <name> -l <login> -p <password> -i <info>",
		Short: "Create a new credential",
		Long:  `Create a new credential with the given name, login, password and optional notes.`,
		RunE:  c.run,
	}

	cmd.Flags().StringVarP(&creds.name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&creds.login, "login", "l", "", "login")
	cmd.Flags().StringVarP(&creds.password, "password", "p", "", "password")
	cmd.Flags().StringVarP(&creds.notes, "info", "i", "", "Info")

	if err := cmd.MarkFlagRequired("name"); err != nil {
		return nil, fmt.Errorf("cli: failed to mark name flag as required: %w", err)
	}
	if err := cmd.MarkFlagRequired("login"); err != nil {
		return nil, fmt.Errorf("cli: failed to mark login flag as required: %w", err)
	}
	if err := cmd.MarkFlagRequired("password"); err != nil {
		return nil, fmt.Errorf("cli: failed to mark password flag as required: %w", err)
	}

	c.cmd = cmd

	return c, nil
}

func (c *CreateCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *CreateCommand) run(cmd *cobra.Command, _ []string) error {
	id := uuid.New().String()
	cred, err := credential.NewCredential(id, creds.name, creds.login, creds.password, creds.notes)
	if err != nil {
		return fmt.Errorf("cli: failed to create credential: %w", err)
	}

	if err := c.credentialSrv.Create(cmd.Context(), *cred); err != nil {
		return fmt.Errorf("cli: failed to create credential: %w", err)
	}

	fmt.Println("Credential created successfully")
	return nil
}
