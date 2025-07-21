package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"buf.build/gen/go/krelinga/proto/connectrpc/go/krelinga/video/in/v1/inv1connect"
	v1 "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestEndToEnd(t *testing.T) {
	ctx := context.Background()

	// Build and start the container using the Dockerfile from parent directory
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "..",
				Dockerfile: "Dockerfile",
			},
			ExposedPorts: []string{"8080/tcp"},
			WaitingFor:   wait.ForListeningPort("8080/tcp").WithStartupTimeout(30 * time.Second),
		},
		Started: true,
	})
	if err != nil {
		t.Fatalf("Failed to start container: %v", err)
	}
	defer func() {
		if err := container.Terminate(ctx); err != nil {
			t.Logf("Failed to terminate container: %v", err)
		}
	}()

	client := func() inv1connect.ServiceClient {
		// Get the mapped port
		mappedPort, err := container.MappedPort(ctx, "8080")
		if err != nil {
			t.Fatalf("Failed to get mapped port: %v", err)
		}

		// Get the container host
		host, err := container.Host(ctx)
		if err != nil {
			t.Fatalf("Failed to get container host: %v", err)
		}

		// Create the service URL
		serviceURL := fmt.Sprintf("http://%s:%s", host, mappedPort.Port())

		// Create a ConnectRPC client
		return inv1connect.NewServiceClient(
			http.DefaultClient,
			serviceURL,
		)
	}()

	t.Run("HelloWorld", func(t *testing.T) {
		// Call the HelloWorld method
		req2 := connect.NewRequest(&v1.HelloWorldRequest{})
		resp, err := client.HelloWorld(ctx, req2)
		if err != nil {
			t.Fatalf("HelloWorld call failed: %v", err)
		}

		// Verify we got a response (even if it's empty as expected from the stub)
		if resp == nil {
			t.Fatal("Expected non-nil response")
		}

		// Verify the response message is not nil
		if resp.Msg == nil {
			t.Fatal("Expected non-nil response message")
		}
	})

	expectedProjects := []string{"Empty", "Name With Spaces"}
	t.Run("ProjectList", func(t *testing.T) {
		// Call the ProjectList method
		req := connect.NewRequest(&v1.ProjectListRequest{})
		resp, err := client.ProjectList(ctx, req)
		if err != nil {
			t.Fatalf("ProjectList call failed: %v", err)
		}

		// Verify we got a response (even if it's empty as expected from the stub)
		if resp == nil {
			t.Fatal("Expected non-nil response")
		}

		// Verify the response message is not nil
		if resp.Msg == nil {
			t.Fatal("Expected non-nil response message")
		}
		for _, ep := range expectedProjects {
			found := false
			for _, p := range resp.Msg.Projects {
				if p == ep {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("Expected project '%s' not found in response", ep)
			}
		}
	})

	t.Run("ProjectGet", func(t *testing.T) {
		for _, projectName := range expectedProjects {
			t.Run(projectName, func(t *testing.T) {
				// Call the ProjectGet method
				req := connect.NewRequest(&v1.ProjectGetRequest{Project: projectName})
				resp, err := client.ProjectGet(ctx, req)
				if err != nil {
					t.Fatalf("ProjectGet call failed for project '%s': %v", projectName, err)
				}

				// Verify we got a response (even if it's empty as expected from the stub)
				if resp == nil {
					t.Fatal("Expected non-nil response")
				}

				// Verify the response message is not nil
				if resp.Msg == nil {
					t.Fatal("Expected non-nil response message")
				}

				// Verify the project name matches
				if resp.Msg.Project != projectName {
					t.Fatalf("Expected project '%s', got '%s'", projectName, resp.Msg.Project)
				}
			})
		}
	})

	t.Run("Check RPC Call in Logs", func(t *testing.T) {
		// Fetch container logs
		logs, err := container.Logs(ctx)
		if err != nil {
			t.Fatalf("Failed to get container logs: %v", err)
		}
		defer logs.Close()

		// Read logs into a string
		logBytes, err := io.ReadAll(logs)
		if err != nil {
			t.Fatalf("Failed to read container logs: %v", err)
		}
		logStr := string(logBytes)

		// Check for 'RPC Call' in logs
		if !strings.Contains(logStr, "RPC Call") {
			t.Fatalf("Expected 'RPC Call' in container logs, but it was not found. Logs:\n%s", logStr)
		}
	})
}
