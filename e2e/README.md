# End-to-End Tests

This module contains end-to-end tests for the video-in-be-stub service using testcontainers.

## Running the Tests

To run the end-to-end tests:

```bash
cd e2e
go test -v .
```

## What the Tests Do

The `TestEndToEnd_HelloWorld` test:

1. Builds the Docker image from the parent directory's Dockerfile
2. Starts a container from the built image
3. Waits for the service to be ready on port 8080
4. Creates a ConnectRPC client to connect to the containerized service
5. Calls the `HelloWorld` method and verifies the response
6. Properly cleans up containers after testing

This validates the complete integration from Docker build through service communication.

## Dependencies

This module uses testcontainers-go to manage Docker containers during testing, which is kept separate from the main module to maintain light dependencies in the core service.