package main

import (
	"os"

	"github.com/nicklasfrahm-dev/appkit/logging"
	"github.com/nicklasfrahm/cloud/cmd/cloud/service"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var help bool

func RootCommand(logger *zap.Logger) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "cloud [command]",
		Short: "cloud manages services in my infrastructure.",
		Long:  `cloud is a command line tool to manage services in my infrastructure.`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if help {
				return cmd.Help()
			}

			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	// Add global flags.
	cmd.PersistentFlags().BoolVarP(&help, "help", "h", false, "Show help for command")

	// Add subcommands.
	cmd.AddCommand(service.RootCommand(logger))

	return cmd
}

func main() {
	// Default to console logger.
	format := os.Getenv("LOG_FORMAT")
	if format == "" {
		os.Setenv("LOG_FORMAT", "console")
	}

	logger := logging.NewLogger()

	err := RootCommand(logger).Execute()
	if err != nil {
		logger.Fatal("Failed to execute command", zap.Error(err))
	}
}
