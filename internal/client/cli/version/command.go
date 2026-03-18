package version

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/version"
	"github.com/spf13/cobra"
)

func NewCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version of the application",
		Long:  `Print the version of the application.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.Info())
		},
	}
}
