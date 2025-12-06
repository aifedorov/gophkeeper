package commands

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewListCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "List all available commands",
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
