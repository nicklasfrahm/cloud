package service

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/nicklasfrahm/cloud/cmd/cloud/workflow"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	deployDir     = "deploy"
	defaultValues = `---
`
)

func Bootstrap(logger *zap.Logger) *cobra.Command {
	var target string

	cmd := &cobra.Command{
		Use:   "bootstrap <service>",
		Short: "Bootstrap the service",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// This should never happen,
			// but better safe than sorry.
			if len(args) != 1 {
				return cmd.Help()
			}

			serviceName := args[0]

			job := workflow.NewJob(
				workflow.EnsureDirectory(path.Join(deployDir, "services", serviceName, "clusters")),
				workflow.EnsureDirectory(path.Join(deployDir, "services", serviceName, "envs")),
				workflow.EnsureFile(
					path.Join(deployDir, "services", serviceName, "envs", "base.yaml"),
					[]byte(defaultValues),
				),
			)

			environments, clusters, err := getEnvironmentsAndClusters()
			if err != nil {
				return fmt.Errorf("failed to get environments and clusters: %w", err)
			}

			for _, environment := range environments {
				job.AddStep(workflow.EnsureFile(
					path.Join(deployDir, "services", serviceName, "envs", environment+".yaml"),
					[]byte(defaultValues),
				))
			}

			for _, cluster := range clusters {
				job.AddStep(workflow.EnsureFile(
					path.Join(deployDir, "services", serviceName, "clusters", cluster+".yaml"),
					[]byte(defaultValues),
				))
			}

			if target != "" {
				chunks := strings.Split(target, "/")
				if len(chunks) != 2 {
					return fmt.Errorf("found invalid target format: expected <environment>/<cluster> but got: %s", target)
				}

				cluster := chunks[0]
				tenant := chunks[1]

				job.AddStep(workflow.EnsureDirectory(
					path.Join(deployDir, "clusters", cluster, tenant, serviceName),
				))
				job.AddStep(workflow.EnsureFile(
					path.Join(deployDir, "clusters", cluster, tenant, serviceName, "config.yaml"),
					[]byte(newDefaultConfig(serviceName)),
				))
				job.AddStep(workflow.EnsureFile(
					path.Join(deployDir, "clusters", cluster, tenant, serviceName, "values.yaml"),
					[]byte(defaultValues),
				))
			}

			err = job.Execute()
			if err != nil {
				return fmt.Errorf("failed to bootstrap service: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "t", "", "The cluster to bootstrap the service in")

	return cmd
}

// getEnvironmentsAndClusters returns a list of environments and clusters.
func getEnvironmentsAndClusters() ([]string, []string, error) {
	var environments []string
	var clusters []string

	entries, err := os.ReadDir(path.Join(deployDir, "clusters"))
	if err != nil {
		return nil, nil, fmt.Errorf("failed to read cluster directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			cluster := entry.Name()
			parts := strings.SplitN(cluster, "-", 2)

			if len(parts) != 2 {
				return nil, nil, fmt.Errorf("found invalid cluster name: %s", cluster)
			}

			environment := parts[0]

			environments = append(environments, environment)
			clusters = append(clusters, cluster)
		}
	}

	return environments, clusters, nil
}

func newDefaultConfig(serviceName string) string {
	return fmt.Sprintf(`chart:
  repo: ghcr.io/nicklasfrahm/charts
  name: %s
  tag: 0.1.0
`, serviceName)
}
