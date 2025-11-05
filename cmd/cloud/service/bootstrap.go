package service

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/nicklasfrahm-dev/platform/cmd/cloud/workflow"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

const (
	deployDir = "deploy"
	// defaultValues is just an empty YAML file with a trailing newline.
	defaultValues = ``
)

func Bootstrap(logger *zap.Logger) *cobra.Command {
	var target string
	var chart string

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
			serviceDir := path.Join(deployDir, "services", serviceName)

			job := workflow.NewJob(
				workflow.EnsureDirectory(serviceDir),
				workflow.EnsureFile(path.Join(serviceDir, "00-base.yml"), []byte(defaultValues)),
			)

			environments, clusters, err := getEnvironmentsAndClusters()
			if err != nil {
				return fmt.Errorf("failed to get environments and clusters: %w", err)
			}

			for _, environment := range environments {
				job.AddStep(workflow.EnsureFile(path.Join(serviceDir, "10-env-"+environment+".yml"), []byte(defaultValues)))
			}

			for _, cluster := range clusters {
				job.AddStep(workflow.EnsureFile(path.Join(serviceDir, "20-cluster-"+cluster+".yml"), []byte(defaultValues)))
			}

			if target != "" {
				chunks := strings.Split(target, "/")
				if len(chunks) != 2 {
					return fmt.Errorf("found invalid target format: expected <cluster>/<tenant> but got: %s", target)
				}

				cluster := chunks[0]
				tenant := chunks[1]
				tenantDir := path.Join(deployDir, "clusters", cluster, tenant)

				job.AddStep(workflow.EnsureDirectory(tenantDir))

				// Default to service name if the chart is not provided.
				if chart == "" {
					chart = serviceName
				}

				job.AddStep(workflow.EnsureFile(path.Join(tenantDir, serviceName+".yml"), []byte(newDefaultConfig(chart))))
			}

			err = job.Execute()
			if err != nil {
				return fmt.Errorf("failed to bootstrap service: %w", err)
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&target, "target", "T", "", "The combination of <cluster>/<tenant> to bootstrap the service in")
	cmd.Flags().StringVarP(&chart, "chart", "C", "", "The chart name to use when deploying the service")

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

func newDefaultConfig(chart string) string {
	return fmt.Sprintf(`chart: %s
tag: 0.1.0
`, chart)
}
