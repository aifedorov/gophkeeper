package text

import (
	"fmt"

	domaintext "github.com/aifedorov/gophkeeper/internal/client/domain/text"
	"github.com/spf13/cobra"
)

type Command struct {
	cmd     *cobra.Command
	textSrv domaintext.Service
}

func NewCommand(textSrv domaintext.Service) (*Command, error) {
	c := &Command{
		textSrv: textSrv,
	}

	cmd := &cobra.Command{
		Use:   "text",
		Short: "Manage text notes",
		Long:  `Manage text notes: create, upload, view, list, update, update-file, delete, download`,
	}

	createCommand, err := NewCreateCommand(textSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create create command: %w", err)
	}
	cmd.AddCommand(createCommand.GetCommand())

	uploadCommand, err := NewUploadCommand(textSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload command: %w", err)
	}
	cmd.AddCommand(uploadCommand.GetCommand())

	viewCommand, err := NewViewCommand(textSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create view command: %w", err)
	}
	cmd.AddCommand(viewCommand.GetCommand())

	listCommand, err := NewListCommand(textSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create list command: %w", err)
	}
	cmd.AddCommand(listCommand.GetCommand())

	updateCommand, err := NewUpdateCommand(textSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create update command: %w", err)
	}
	cmd.AddCommand(updateCommand.GetCommand())

	updateFileCommand, err := NewUpdateFileCommand(textSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create update-file command: %w", err)
	}
	cmd.AddCommand(updateFileCommand.GetCommand())

	deleteCommand, err := NewDeleteCommand(textSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create delete command: %w", err)
	}
	cmd.AddCommand(deleteCommand.GetCommand())

	downloadCommand, err := NewDownloadCommand(textSrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create download command: %w", err)
	}
	cmd.AddCommand(downloadCommand.GetCommand())

	c.cmd = cmd

	return c, nil
}

func (c *Command) GetCommand() *cobra.Command {
	return c.cmd
}
