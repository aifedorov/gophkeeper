package credential

import (
	"fmt"
	"strings"

	"github.com/aifedorov/gophkeeper/internal/client/cli/shared"
	"github.com/aifedorov/gophkeeper/internal/client/domain/credential"
	"github.com/spf13/cobra"
)

type ListCommand struct {
	cmd           *cobra.Command
	credentialSrv credential.Service
}

func NewListCommand(credentialSrv credential.Service) (*ListCommand, error) {
	c := &ListCommand{
		credentialSrv: credentialSrv,
	}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all credentials",
		Long:  `List all credentials`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd)
		},
	}

	c.cmd = cmd

	return c, nil
}

func (c *ListCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *ListCommand) run(cmd *cobra.Command) error {
	creds, err := c.credentialSrv.List(cmd.Context())
	if err != nil {
		return shared.ParseErrorForCLI(err)
	}

	if len(creds) == 0 {
		fmt.Println("No credentials found")
		return nil
	}

	fmt.Printf("%-36s %-20s %-30s %-15s %s\n", "ID", "NAME", "LOGIN", "PASSWORD", "NOTES")
	fmt.Println(strings.Repeat("-", 120))

	for _, cred := range creds {
		fmt.Printf("%-36s %-20s %-30s %-15s %s\n", cred.ID, cred.Name, cred.Login, cred.Password, cred.Notes)
	}

	return nil
}
