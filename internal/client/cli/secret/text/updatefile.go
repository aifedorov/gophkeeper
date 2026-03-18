package text

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/text"
	"github.com/spf13/cobra"
)

type UpdateFileCommand struct {
	cmd     *cobra.Command
	textSrv text.Service
}

func NewUpdateFileCommand(textSrv text.Service) (*UpdateFileCommand, error) {
	c := &UpdateFileCommand{
		textSrv: textSrv,
	}

	var id string
	var filePath string
	var notes string

	cmd := &cobra.Command{
		Use:   "update-file -i <id> -f <file> [-n <notes>]",
		Short: "Update an existing text note from a file",
		Long:  `Update an existing text note by uploading a new file.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, id, filePath, notes)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID of the text note to update (required)")
	err := cmd.MarkFlagRequired("id")
	if err != nil {
		return nil, fmt.Errorf("`id` flag as required: %w", err)
	}

	cmd.Flags().StringVarP(&filePath, "file", "f", "", "Path to file to upload (required)")
	err = cmd.MarkFlagRequired("file")
	if err != nil {
		return nil, fmt.Errorf("`file` flag as required: %w", err)
	}

	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Optional notes about the text note")

	c.cmd = cmd

	return c, nil
}

func (c *UpdateFileCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *UpdateFileCommand) run(cmd *cobra.Command, id, filePath, notes string) error {
	if id == "" {
		return fmt.Errorf("cli: id is required")
	}
	if filePath == "" {
		return fmt.Errorf("cli: file path is required")
	}

	fmt.Printf("Updating text note '%s' from file '%s'...\n", id, filePath)
	if err := c.textSrv.UpdateFromFile(cmd.Context(), id, filePath, notes); err != nil {
		return fmt.Errorf("cli: failed to update text note: %w", err)
	}

	fmt.Println("Text note updated successfully")
	return nil
}
