package text

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/text"
	"github.com/spf13/cobra"
)

type DownloadCommand struct {
	cmd     *cobra.Command
	textSrv text.Service
}

func NewDownloadCommand(textSrv text.Service) (*DownloadCommand, error) {
	c := &DownloadCommand{
		textSrv: textSrv,
	}

	var id string

	cmd := &cobra.Command{
		Use:   "download -i <id>",
		Short: "Download a text note",
		Long:  `Download a text note by ID and save it to local storage.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, id)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID of the text note to download (required)")
	err := cmd.MarkFlagRequired("id")
	if err != nil {
		return nil, fmt.Errorf("`id` flag as required: %w", err)
	}

	c.cmd = cmd

	return c, nil
}

func (c *DownloadCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *DownloadCommand) run(cmd *cobra.Command, id string) error {
	if id == "" {
		return fmt.Errorf("cli: id is required")
	}

	filepath, err := c.textSrv.Download(cmd.Context(), id)
	if err != nil {
		return fmt.Errorf("cli: failed to download text note: %w", err)
	}

	fmt.Printf("Text note downloaded successfully to: %s\n", filepath)

	return nil
}
