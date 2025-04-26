// Package server contains the plumbing for a server
// that can handle both gRPC and REST requests.
package server

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/nicklasfrahm/cloud/pkg/kms"
	"github.com/soheilhy/cmux"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const defaultPort = 8080

// ServerListeners constains the listeners for the server.
type ServerListeners struct {
	GRPCListener net.Listener
	HTTPListener net.Listener
}

// Server is a struct that contains the server configuration.
type Server struct {
	listerners *ServerListeners
}

// New creates a new server instance.
func New() *Server {

	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	m := cmux.New(listener)

	grpcL := m.MatchWithWriters(cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc"))
	httpL := m.Match(cmux.Any()) // Now handle both HTTP/1 and h2c in one place

	// Set up REST handlers
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Pong from REST (h2c supported)\n")
	})

	// Wrap the HTTP handler with h2c support
	h2s := &http2.Server{}
	h1h2cHandler := h2c.NewHandler(mux, h2s)

	// gRPC server
	grpcServer := grpc.NewServer()
	taloskms.RegisterKMSServiceServer(grpcServer, &kms.KMSServiceServer{})

	// Enable gRPC reflection
	reflection.Register(grpcServer)

	// Start gRPC
	go func() {
		log.Println("Starting gRPC server on :" + port)
		if err := grpcServer.Serve(grpcL); err != nil {
			log.Fatalf("gRPC server failed: %v", err)
		}
	}()

	// Start REST with h2c support
	go func() {
		log.Println("Starting REST server on :" + port)
		if err := http.Serve(httpL, h1h2cHandler); err != nil {
			log.Fatalf("REST server failed: %v", err)
		}
	}()

	log.Println("cmux multiplexer running on :" + port)
	if err := m.Serve(); err != nil {
		log.Fatalf("cmux server failed: %v", err)
	}
}
