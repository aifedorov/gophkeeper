package text

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/text"
	"github.com/spf13/cobra"
)

type UploadCommand struct {
	cmd     *cobra.Command
	textSrv text.Service
}

func NewUploadCommand(textSrv text.Service) (*UploadCommand, error) {
	c := &UploadCommand{
		textSrv: textSrv,
	}

	var filePath string
	var notes string

	cmd := &cobra.Command{
		Use:   "upload -f <file> [-n <notes>]",
		Short: "Upload a text file",
		Long:  `Upload a text file from the filesystem.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, filePath, notes)
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to file to upload (required)")
	err := cmd.MarkFlagRequired("file")
	if err != nil {
		return nil, fmt.Errorf("`file` flag as required: %w", err)
	}

	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Optional notes about the text file")

	c.cmd = cmd

	return c, nil
}

func (c *UploadCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *UploadCommand) run(cmd *cobra.Command, filePath, notes string) error {
	if filePath == "" {
		return fmt.Errorf("cli: file path is required")
	}

	fmt.Printf("Uploading text file '%s'...\n", filePath)
	if err := c.textSrv.CreateFromFile(cmd.Context(), filePath, notes); err != nil {
		return fmt.Errorf("cli: failed to upload text file: %w", err)
	}

	fmt.Println("Text file uploaded successfully")
	return nil
}
