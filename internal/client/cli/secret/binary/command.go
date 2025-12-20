package binary

import (
	"fmt"

	domainbinary "github.com/aifedorov/gophkeeper/internal/client/domain/binary"
	"github.com/spf13/cobra"
)

type Command struct {
	cmd       *cobra.Command
	binarySrv domainbinary.Service
}

func NewCommand(binarySrv domainbinary.Service) (*Command, error) {
	c := &Command{
		binarySrv: binarySrv,
	}

	cmd := &cobra.Command{
		Use:   "file",
		Short: "Manage files",
		Long:  `Manage files: upload, download, list, delete`,
	}

	uploadCommand, err := NewUploadCommand(binarySrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload command: %w", err)
	}
	cmd.AddCommand(uploadCommand.GetCommand())

	listCommand, err := NewListCommand(binarySrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create list command: %w", err)
	}
	cmd.AddCommand(listCommand.GetCommand())

	downloadCommand, err := NewDownloadCommand(binarySrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create download command: %w", err)
	}
	cmd.AddCommand(downloadCommand.GetCommand())

	c.cmd = cmd

	deleteCommand, err := NewDeleteCommand(binarySrv)
	if err != nil {
		return nil, fmt.Errorf("failed to create delete command: %w", err)
	}
	cmd.AddCommand(deleteCommand.GetCommand())

	return c, nil
}

func (c *Command) GetCommand() *cobra.Command {
	return c.cmd
}
