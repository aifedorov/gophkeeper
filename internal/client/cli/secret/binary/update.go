package binary

import (
	"fmt"
	"os"

	"github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/spf13/cobra"
)

type UpdateCommand struct {
	cmd       *cobra.Command
	binarySrv binary.Service
}

func NewUpdateCommand(binarySrv binary.Service) (*UpdateCommand, error) {
	c := &UpdateCommand{
		binarySrv: binarySrv,
	}

	var id string
	var filePath string
	var notes string

	cmd := &cobra.Command{
		Use:   "update -i <id> -f <file> [-n <notes>]",
		Short: "Update a file",
		Long:  `Update an existing file on the server with a new version.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, id, filePath, notes)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "File ID to update (required)")
	err := cmd.MarkFlagRequired("id")
	if err != nil {
		return nil, fmt.Errorf("`id` flag as required: %w", err)
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to new file (required)")
	err = cmd.MarkFlagRequired("file")
	if err != nil {
		return nil, fmt.Errorf("`file` flag as required: %w", err)
	}

	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Optional notes about the file")

	c.cmd = cmd

	return c, nil
}

func (c *UpdateCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *UpdateCommand) run(cmd *cobra.Command, id string, filePath string, notes string) error {
	if id == "" {
		return fmt.Errorf("cli: file ID is required")
	}
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

	fmt.Printf("Updating file %s with: %s (size: %d bytes)...\n", id, stat.Name(), stat.Size())
	if err := c.binarySrv.Update(cmd.Context(), id, filePath, notes); err != nil {
		return fmt.Errorf("cli: failed to update file: %w", err)
	}

	fmt.Println("File updated successfully")

	return nil
}
