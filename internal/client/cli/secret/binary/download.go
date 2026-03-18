package binary

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/spf13/cobra"
)

type DownloadCommand struct {
	cmd       *cobra.Command
	binarySrv binary.Service
}

func NewDownloadCommand(binarySrv binary.Service) (*DownloadCommand, error) {
	c := &DownloadCommand{
		binarySrv: binarySrv,
	}

	var fileID string

	cmd := &cobra.Command{
		Use:   "download -i <file-id>",
		Short: "Download a file",
		Long:  `Download a file from the server to local storage.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, fileID)
		},
	}

	cmd.Flags().StringVarP(&fileID, "id", "i", "", "ID of the file to download (required)")
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

func (c *DownloadCommand) run(cmd *cobra.Command, fileID string) error {
	if len(fileID) < 32 {
		return fmt.Errorf("cli: invalid file ID format. Expected UUID")
	}

	fmt.Printf("Downloading file with ID: %s...\n", fileID)
	filepath, err := c.binarySrv.Download(cmd.Context(), fileID)
	if err != nil {
		return fmt.Errorf("cli: failed to download file: %w", err)
	}

	fmt.Printf("File downloaded successfully to: %s\n", filepath)

	return nil
}
