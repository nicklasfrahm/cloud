package main

import (
	_ "github.com/joho/godotenv/autoload"
	"github.com/nicklasfrahm/cloud/pkg/zone"
	"github.com/spf13/cobra"
)

var zoneUpCmd = &cobra.Command{
	Use:   "up <host>",
	Short: "Bootstrap a new availability zone",
	Long: `This command will bootstrap a new zone by connecting
to the specified IP and setting up a k3s cluster on
the host that will then set up the required services
for managing the lifecycle of the zone.
To manage a zone, the CLI needs credentials for the
DNS provider that is used to manage the DNS records
for the zone. These credentials can only be provided
via the environment variable DNS_PROVIDER_CREDENTIAL
and DNS_PROVIDER or via a ".env" file in the current
working directory.`,
	Args:       cobra.ExactArgs(1),
	ArgAliases: []string{"host"},
	ValidArgs:  []string{"host"},
	RunE: func(cmd *cobra.Command, args []string) error {
		return zone.Up(zone.Zone{
			Name: args[0],
		})
	},
}
