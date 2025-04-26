package server

import (
	"context"
	"os"
	"strconv"
)

// DefaultPort is the default port for the server.
const DefaultPort = 8080

// getPort returns the port to listen on. It checks the
// PORT environment variable first, then defaults to 8080.
func getPort(ctx context.Context) int {
	port := os.Getenv("PORT")
	if port == "" {
		return DefaultPort
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {

	}

}
