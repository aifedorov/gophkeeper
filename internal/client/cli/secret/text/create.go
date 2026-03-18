package text

import (
	"fmt"

	"github.com/aifedorov/gophkeeper/internal/client/domain/text"
	"github.com/spf13/cobra"
)

type CreateCommand struct {
	cmd     *cobra.Command
	textSrv text.Service
}

func NewCreateCommand(textSrv text.Service) (*CreateCommand, error) {
	c := &CreateCommand{
		textSrv: textSrv,
	}

	var content string
	var title string
	var notes string

	cmd := &cobra.Command{
		Use:   "create -c <content> -t <title> [-n <notes>]",
		Short: "Create a new text note from inline content",
		Long:  `Create a new text note from inline content.`,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return c.run(cmd, content, title, notes)
		},
	}

	cmd.Flags().StringVarP(&content, "content", "c", "", "Inline content for the note (required)")
	err := cmd.MarkFlagRequired("content")
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

func (c *CreateCommand) GetCommand() *cobra.Command {
	return c.cmd
}

func (c *CreateCommand) run(cmd *cobra.Command, content, title, notes string) error {
	if content == "" {
		return fmt.Errorf("cli: content is required")
	}
	if title == "" {
		return fmt.Errorf("cli: title is required")
	}

	fmt.Printf("Creating text note '%s' from inline content...\n", title)
	if err := c.textSrv.CreateFromContent(cmd.Context(), content, title, notes); err != nil {
		return fmt.Errorf("cli: failed to create text note: %w", err)
	}

	fmt.Println("Text note created successfully")
	return nil
}
