package root

import (
	"github.com/spf13/cobra"
)

type RootCommand struct {
	cmd *cobra.Command
}

func NewCommand() *RootCommand {
	cmd := &cobra.Command{
		Use:   "gophkeeper",
		Short: "GophKeeper is a secure password manager",
		Long:  `GophKeeper is a secure password manager that allows you to store and retrieve your passwords securely.`,
	}
	return &RootCommand{
		cmd: cmd,
	}
}

func (r *RootCommand) AddCommand(cmd *cobra.Command) {
	r.cmd.AddCommand(cmd)
}

func (r *RootCommand) Execute() error {
	return r.cmd.Execute()
}
