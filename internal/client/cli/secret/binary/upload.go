package binary

import (
	"fmt"
	"os"

	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/spf13/cobra"
)

type UploadCommand struct {
	cmd       *cobra.Command
	binarySrv binary.Service
}

func NewUploadCommand(binarySrv binary.Service) (*UploadCommand, error) {
	c := &UploadCommand{
		binarySrv: binarySrv,
	}

	var filePath string

	cmd := &cobra.Command{
		Use:   "upload -f <file>",
		Short: "Upload a file",
		Long:  `Upload a file to the server for secure storage.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, filePath)
		},
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to file to upload (required)")
	_ = cmd.MarkFlagRequired("file")

	c.cmd = cmd

	return c, nil
}

func (c *UploadCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *UploadCommand) run(cmd *cobra.Command, filePath string) error {
	if filePath == "" {
		return fmt.Errorf("cli: file path is required")
	}
	stat, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return fmt.Errorf("cli: file does not exist: %s", filePath)
	}
	if err != nil {
		return fmt.Errorf("cli: failed to get file stat: %w", err)
	}

	fmt.Printf("Uploading file: %s with size: %d bytes...\n", stat.Name(), stat.Size())
	if err := c.binarySrv.Upload(cmd.Context(), filePath); err != nil {
		return fmt.Errorf("cli: failed to upload file: %w", err)
	}

	fmt.Println("File uploaded successfully")

	return nil
}
