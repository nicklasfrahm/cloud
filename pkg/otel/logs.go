package otel

import (
	"context"
	"os"
	"strings"

	"github.com/nicklasfrahm/cloud/cmd/otel/zap"
	"go.opentelemetry.io/otel/sdk/log"
)

// NewLoggerProvider defaults to a console exporter if no exporter is specified.
// If an OTLP endpoint is specified, it will use that endpoint in addtion to the
// console exporter. This is useful to allow developers to still inspect container
// logs in the console while also sending them to a remote endpoint.
func NewLoggerProvider(ctx context.Context) *sdk.LoggerProvider {
	rawExporterNames := os.Getenv("OTLP_LOGS_EXPORTER")
	if rawExporterNames == "" {
		rawExporterNames = "console"
	}

	exporterNames := strings.Split(rawExporterNames, ",")

	otlpEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_LOGS_ENDPOINT")
	if otlpEndpoint == "" {
		otlpEndpoint = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	}

	exporters := make([]log.Exporter, 0, len(exporterNames))
	for _, exporterName := range exporterNames {
		switch strings.ToLower(exporterName) {
		case "console":
			exporters = append(exporters, zap.New())
		case "otlp":
			// TODO: Configure the OTLP exporter.
			exporters = append(exporters, newOtlpExporter(ctx))
		default:
			continue
		}
	}

	if otlpEndpoint != "" {
		exporters = append(exporters, "otlp")
	}

	loggerProvider := sdk.NewLoggerProvider()

	return loggerProvider
}
