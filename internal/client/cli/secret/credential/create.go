package credential

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CreateCommand struct {
	cmd           *cobra.Command
	credentialSrv credential.Service
}

func NewCreateCommand(credentialSrv credential.Service) (*CreateCommand, error) {
	c := &CreateCommand{
		credentialSrv: credentialSrv,
	}

	cred := inputCredentials{}

	cmd := &cobra.Command{
		Use:   "create -n <name> -l <login> -p <password> -i <info>",
		Short: "Create a new credential",
		Long:  `Create a new credential with the given name, login, password and optional notes.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, cred)
		},
	}

	cmd.Flags().StringVarP(&cred.name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&cred.login, "login", "l", "", "login")
	cmd.Flags().StringVarP(&cred.password, "password", "p", "", "password")
	cmd.Flags().StringVarP(&cred.notes, "info", "i", "", "Info")

	c.cmd = cmd

	return c, nil
}

func (c *CreateCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *CreateCommand) run(cmd *cobra.Command, cred inputCredentials) error {
	id := uuid.New().String()
	newCred, err := credential.NewCredential(id, cred.name, cred.login, cred.password, cred.notes)
	if err != nil || newCred == nil {
		return fmt.Errorf("cli: failed to create credential: %w", err)
	}

	if err := newCred.Validate(); err != nil {
		return fmt.Errorf("cli: failed to validate credentials: %w", err)
	}

	if err := c.credentialSrv.Create(cmd.Context(), *newCred); err != nil {
		return fmt.Errorf("cli: failed to create credential: %w", err)
	}

	fmt.Println("Credential created successfully")
	return nil
}
