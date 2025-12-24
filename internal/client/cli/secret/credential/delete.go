package credential

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/cli/shared"
	"github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	"github.com/spf13/cobra"
)

type DeleteCommand struct {
	cmd           *cobra.Command
	credentialSrv credential.Service
}

func NewDeleteCommand(credentialSrv credential.Service) (*DeleteCommand, error) {
	c := &DeleteCommand{
		credentialSrv: credentialSrv,
	}

	var id string
	cmd := &cobra.Command{
		Use:   "delete -i <id>",
		Short: "Delete a credential by ID",
		Long:  `Delete a credential by ID`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, id)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID")
	if err := cmd.MarkFlagRequired("id"); err != nil {
		return nil, fmt.Errorf("cli: failed to mark id flag as required: %w", err)
	}

	c.cmd = cmd

	return c, nil
}

func (c *DeleteCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *DeleteCommand) run(cmd *cobra.Command, id string) error {
	err := c.credentialSrv.Delete(cmd.Context(), id)
	if err != nil {
		return shared.ParseErrorForCLI(err)
	}

	fmt.Println("Credential deleted successfully")
	return nil
}
