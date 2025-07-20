package main

import (
	"bytes"
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	v1 "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"buf.build/gen/go/krelinga/proto/connectrpc/go/krelinga/video/in/v1/inv1connect"
	"connectrpc.com/connect"
)

func TestLoggingInterceptor(t *testing.T) {
	// Capture log output using a buffer as the log output
	var buf bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&buf)
	defer log.SetOutput(originalOutput)

	// Create service with logging interceptor
	service := NewStubService()
	loggingInterceptor := &LoggingInterceptor{}
	
	// Create handler with interceptor
	_, handler := inv1connect.NewServiceHandler(
		service,
		connect.WithInterceptors(loggingInterceptor),
	)

	// Create test server
	server := httptest.NewServer(handler)
	defer server.Close()

	// Create client
	client := inv1connect.NewServiceClient(http.DefaultClient, server.URL)
	ctx := context.Background()

	// Test successful request
	req := connect.NewRequest(&v1.HelloWorldRequest{Name: "test"})
	resp, err := client.HelloWorld(ctx, req)
	
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	
	if resp == nil {
		t.Fatal("Expected non-nil response")
	}

	// Verify logging occurred
	logOutput := buf.String()
	if !strings.Contains(logOutput, "RPC Call") {
		t.Errorf("Expected log output to contain 'RPC Call', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "HelloWorld") {
		t.Errorf("Expected log output to contain method name 'HelloWorld', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "test") {
		t.Errorf("Expected log output to contain request data 'test', got: %s", logOutput)
	}
	if !strings.Contains(logOutput, "Hello, test!") {
		t.Errorf("Expected log output to contain response data 'Hello, test!', got: %s", logOutput)
	}

	// Reset buffer for error test
	buf.Reset()

	// Test error case
	errorReq := connect.NewRequest(&v1.HelloWorldRequest{Name: "unknown"})
	_, err = client.HelloWorld(ctx, errorReq)
	
	if err == nil {
		t.Fatal("Expected error for unknown request")
	}

	// Verify error logging
	errorLogOutput := buf.String()
	if !strings.Contains(errorLogOutput, "RPC Call") {
		t.Errorf("Expected error log output to contain 'RPC Call', got: %s", errorLogOutput)
	}
	if !strings.Contains(errorLogOutput, "Error") {
		t.Errorf("Expected error log output to contain 'Error', got: %s", errorLogOutput)
	}
}