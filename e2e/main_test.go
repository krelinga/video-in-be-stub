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

	expectedUnclaimed := []string{"Unclaimed1", "Unclaimed 2"}
	t.Run("UnclaimedDiscDirList", func(t *testing.T) {
		// Call the UnclaimedDiscDirList method
		req := connect.NewRequest(&v1.UnclaimedDiscDirListRequest{})
		resp, err := client.UnclaimedDiscDirList(ctx, req)
		if err != nil {
			t.Fatalf("UnclaimedDiscDirList call failed: %v", err)
		}
		// Verify we got a response (even if it's empty as expected from the stub)
		if resp == nil {
			t.Fatal("Expected non-nil response")
		}
		// Verify the response message is not nil
		if resp.Msg == nil {
			t.Fatal("Expected non-nil response message")
		}
		for _, eu := range expectedUnclaimed {
			found := false
			for _, u := range resp.Msg.Dirs {
				if u == eu {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("Expected unclaimed '%s' not found in response", eu)
			}
		}
	})

	expectedMovieTitles := []string{"Movie 1", "Movie 2"}
	t.Run("MovieSearch", func(t *testing.T) {
		// Call the MovieSearch method with a partial title
		req := connect.NewRequest(&v1.MovieSearchRequest{PartialTitle: "Movie"})
		resp, err := client.MovieSearch(ctx, req)
		if err != nil {
			t.Fatalf("MovieSearch call failed: %v", err)
		}

		// Verify we got a response (even if it's empty as expected from the stub)
		if resp == nil {
			t.Fatal("Expected non-nil response")
		}

		// Verify the response message is not nil
		if resp.Msg == nil {
			t.Fatal("Expected non-nil response message")
		}

		for _, title := range expectedMovieTitles {
			found := false
			for _, result := range resp.Msg.Results {
				if result.Title == title {
					found = true
					break
				}
			}
			if !found {
				t.Fatalf("Expected movie title '%s' not found in response", title)
			}
		}
	})

	t.Run("UnimplementedMethods", func(t *testing.T) {
		calls := []struct {
			method string
			call   func() error
		}{
			{
				method: "ProjectNew",
				call: func() error {
					req := connect.NewRequest(&v1.ProjectNewRequest{})
					_, err := client.ProjectNew(ctx, req)
					return err
				},
			},
			{
				method: "ProjectAssignDiskDirs",
				call: func() error {
					req := connect.NewRequest(&v1.ProjectAssignDiskDirsRequest{})
					_, err := client.ProjectAssignDiskDirs(ctx, req)
					return err
				},
			},
			{
				method: "ProjectCategorizeFiles",
				call: func() error {
					req := connect.NewRequest(&v1.ProjectCategorizeFilesRequest{})
					_, err := client.ProjectCategorizeFiles(ctx, req)
					return err
				},
			},
			{
				method: "ProjectSetMetadata",
				call: func() error {
					req := connect.NewRequest(&v1.ProjectSetMetadataRequest{})
					_, err := client.ProjectSetMetadata(ctx, req)
					return err
				},
			},
			{
				method: "ProjectFinish",
				call: func() error {
					req := connect.NewRequest(&v1.ProjectFinishRequest{})
					_, err := client.ProjectFinish(ctx, req)
					return err
				},
			},
			{
				method: "ProjectAbandon",
				call: func() error {
					req := connect.NewRequest(&v1.ProjectAbandonRequest{})
					_, err := client.ProjectAbandon(ctx, req)
					return err
				},
			},
		}
		for _, c := range calls {
			t.Run(c.method, func(t *testing.T) {
				err := c.call()
				if err == nil {
					t.Fatalf("Expected error for unimplemented method '%s', got nil", c.method)
				}
				if !strings.Contains(err.Error(), "is not implemented") {
					t.Fatalf("Expected unimplemented error for method '%s', got: %v", c.method, err)
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
