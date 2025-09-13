package service

import (
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func RootCommand(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "service [command]",
		Short: "Manage the lifecycle of a service",
		Long:  `Manage the lifecycle of a service.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(Bootstrap(logger))

	return cmd
}
