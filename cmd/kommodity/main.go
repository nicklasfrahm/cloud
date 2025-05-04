package main

import (
	"context"
	"net/http"

	"github.com/nicklasfrahm/cloud/pkg/kms"
	"github.com/nicklasfrahm/cloud/pkg/otel"
	"github.com/nicklasfrahm/cloud/pkg/server"
	taloskms "github.com/siderolabs/kms-client/api/kms"
	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func main() {
	ctx := context.Background()

	// Configure opentelemetry logger provider.
	loggerProvider := otel.NewLoggerProvider(ctx)
	defer loggerProvider.Shutdown(ctx)

	logger := zap.New(otelzap.NewCore("kommodity", otelzap.WithLoggerProvider(loggerProvider)))
	zap.ReplaceGlobals(logger)

	srv := NewServer(ctx)

	if err := srv.ListenAndServe(ctx); err != nil {
		logger.Fatal("failed to start server", zap.Error(err))
	}
}

func NewServer(ctx context.Context) *server.Server {
	srv := server.New(ctx).
		WithGRPCServerInitializer(func(grpcServer *grpc.Server) error {
			taloskms.RegisterKMSServiceServer(grpcServer, &kms.KMSServiceServer{})

			return nil
		}).
		WithHTTPMuxInitializer(func(mux *http.ServeMux) error {
			mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
				w.Write([]byte("OK"))
			})

			return nil
		})

	return srv
}
