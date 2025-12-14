package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewAllCommandsCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "commands",
		Short: "List all available commands with their descriptions",
		Long:  `List all available commands with their descriptions.`,
		Run: func(cmd *cobra.Command, args []string) {
			commands := cmd.Root().Commands()
			fmt.Println("Available commands:")
			for _, c := range commands {
				fmt.Printf("  %-15s %s\n", c.Name(), c.Short)
			}
		},
	}
}
