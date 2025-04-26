package main

import (
	"context"
	"log"

	"github.com/nicklasfrahm/cloud/pkg/server"
)

func main() {
	ctx := context.Background()

	// configure opentelemetry logger provider
	logExporter, _ := otlplogs.NewExporter(ctx)
	loggerProvider := sdk.NewLoggerProvider(
		sdk.WithBatcher(logExporter),
		sdk.WithResource(newResource()),
	)
	// gracefully shutdown logger to flush accumulated signals before program finish
	defer loggerProvider.Shutdown(ctx)

	srv := server.New()

	if err := srv.ListenAndServe(ctx); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
