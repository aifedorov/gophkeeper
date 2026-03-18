package text

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/text"
	"github.com/spf13/cobra"
)

type UpdateCommand struct {
	cmd     *cobra.Command
	textSrv text.Service
}

func NewUpdateCommand(textSrv text.Service) (*UpdateCommand, error) {
	c := &UpdateCommand{
		textSrv: textSrv,
	}

	var id string
	var content string
	var title string
	var notes string

	cmd := &cobra.Command{
		Use:   "update -i <id> -c <content> -t <title> [-n <notes>]",
		Short: "Update an existing text note with inline content",
		Long:  `Update an existing text note with inline content.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, id, content, title, notes)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "ID of the text note to update (required)")
	err := cmd.MarkFlagRequired("id")
	if err != nil {
		return nil, fmt.Errorf("`id` flag as required: %w", err)
	}

	cmd.Flags().StringVarP(&content, "content", "c", "", "Inline content for the note (required)")
	err = cmd.MarkFlagRequired("content")
	if err != nil {
		return nil, fmt.Errorf("`content` flag as required: %w", err)
	}

	cmd.Flags().StringVarP(&title, "title", "t", "", "Title/name for the note (required)")
	err = cmd.MarkFlagRequired("title")
	if err != nil {
		return nil, fmt.Errorf("`title` flag as required: %w", err)
	}

	cmd.Flags().StringVarP(&notes, "notes", "n", "", "Optional notes about the text note")

	c.cmd = cmd

	return c, nil
}

func (c *UpdateCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *UpdateCommand) run(cmd *cobra.Command, id, content, title, notes string) error {
	if id == "" {
		return fmt.Errorf("cli: id is required")
	}
	if content == "" {
		return fmt.Errorf("cli: content is required")
	}
	if title == "" {
		return fmt.Errorf("cli: title is required")
	}

	fmt.Printf("Updating text note '%s' with inline content...\n", id)
	if err := c.textSrv.UpdateFromContent(cmd.Context(), id, content, title, notes); err != nil {
		return fmt.Errorf("cli: failed to update text note: %w", err)
	}

	fmt.Println("Text note updated successfully")
	return nil
}
