package main_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	main "github.com/nicklasfrahm/cloud/cmd/kommodity"
	taloskms "github.com/siderolabs/kms-client/api/kms"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestKMSService(t *testing.T) {
	// Arrange.
	port := "50051"
	t.Setenv("PORT", port)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Start the server in a goroutine
	go func() {
		main.NewServer(ctx).ListenAndServe(ctx)
	}()

	// Query the /health endpoint until we get a 200 response.
	retries := 0
	for range retries {
		if retries > 10 {
			t.Fatalf("Server failed to start")
		}

		resp, err := http.Get("http://localhost:" + port + "/health")
		require.NoError(t, err, "Failed to query /health endpoint")

		if resp.StatusCode == http.StatusOK {
			break
		}

		time.Sleep(100 * time.Millisecond)
	}

	// Create a client connection.
	conn, err := grpc.NewClient("localhost:"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err, "Failed to connect to server")

	t.Cleanup(func() {
		conn.Close()
	})

	// Create a client.
	client := taloskms.NewKMSServiceClient(conn)

	// Act: Test Seal
	sealReq := &taloskms.Request{
		Data: []byte("test data"),
	}
	sealResp, err := client.Seal(ctx, sealReq)

	// Assert: Seal.
	require.NoError(t, err, "Seal failed")
	assert.Equal(t, "sealed:test data", string(sealResp.Data), "Unexpected seal response")

	// Act: Test Unseal.
	unsealReq := &taloskms.Request{
		Data: sealResp.Data,
	}
	unsealResp, err := client.Unseal(ctx, unsealReq)

	// Assert: Unseal.
	require.NoError(t, err, "Unseal failed")
	assert.Equal(t, "test data", string(unsealResp.Data), "Unexpected unseal response")
}
