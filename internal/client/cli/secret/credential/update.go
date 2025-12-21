package credential

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	"github.com/spf13/cobra"
)

type UpdateCommand struct {
	cmd           *cobra.Command
	credentialSrv credential.Service
}

func NewUpdateCommand(credentialSrv credential.Service) (*UpdateCommand, error) {
	c := &UpdateCommand{
		credentialSrv: credentialSrv,
	}

	cred := inputCredentials{}

	cmd := &cobra.Command{
		Use:   "update -i <id> -n <name> -l <login> -p <password> -f <info>",
		Short: "Update a credential by ID",
		Long:  `Update a credential by ID with the given name, login, password and optional notes.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, cred)
		},
	}

	cmd.Flags().StringVarP(&cred.id, "id", "i", "", "ID")
	cmd.Flags().StringVarP(&cred.name, "name", "n", "", "Name")
	cmd.Flags().StringVarP(&cred.login, "login", "l", "", "login")
	cmd.Flags().StringVarP(&cred.password, "password", "p", "", "password")
	cmd.Flags().StringVarP(&cred.notes, "info", "f", "", "Info")

	c.cmd = cmd

	return c, nil
}

func (c *UpdateCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *UpdateCommand) run(cmd *cobra.Command, cred inputCredentials) error {
	if err := cred.Validate(); err != nil {
		return fmt.Errorf("cli: failed to validate credentials: %w", err)
	}

	err := c.credentialSrv.Update(cmd.Context(), cred.id, credential.Credential{
		ID:       cred.id,
		Name:     cred.name,
		Login:    cred.login,
		Password: cred.password,
		Notes:    cred.notes,
	})
	if err != nil {
		return fmt.Errorf("cli: failed to update credential: %w", err)
	}

	fmt.Println("Credential updated successfully")
	return nil
}
