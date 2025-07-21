package main

import (
	"context"
	"errors"
	"testing"

	v1 "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"google.golang.org/protobuf/proto"
)

func TestStubService_HelloWorld_Mapping(t *testing.T) {
	service := NewStubService()
	ctx := context.Background()

	tests := []struct {
		name           string
		request        *v1.HelloWorldRequest
		expectedResp   *v1.HelloWorldResponse
		expectError    bool
		expectedErrCode connect.Code
	}{
		{
			name:         "empty request matches mapping",
			request:      &v1.HelloWorldRequest{Name: ""},
			expectedResp: &v1.HelloWorldResponse{Message: "Hello, empty!"},
			expectError:  false,
		},
		{
			name:         "test request matches mapping",
			request:      &v1.HelloWorldRequest{Name: "test"},
			expectedResp: &v1.HelloWorldResponse{Message: "Hello, test!"},
			expectError:  false,
		},
		{
			name:         "world request matches mapping",
			request:      &v1.HelloWorldRequest{Name: "world"},
			expectedResp: &v1.HelloWorldResponse{Message: "Hello, world!"},
			expectError:  false,
		},
		{
			name:            "unknown request returns NOT_FOUND",
			request:         &v1.HelloWorldRequest{Name: "unknown"},
			expectError:     true,
			expectedErrCode: connect.CodeNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := connect.NewRequest(tt.request)
			resp, err := service.HelloWorld(ctx, req)

			if tt.expectError {
				if err == nil {
					t.Fatal("Expected error but got none")
				}
				var connectErr *connect.Error
				if !errors.As(err, &connectErr) {
					t.Fatal("Expected connect.Error but got different error type")
				}
				if connectErr.Code() != tt.expectedErrCode {
					t.Fatalf("Expected error code %v, got %v", tt.expectedErrCode, connectErr.Code())
				}
			} else {
				if err != nil {
					t.Fatalf("Unexpected error: %v", err)
				}
				if resp == nil {
					t.Fatal("Expected non-nil response")
				}
				if !proto.Equal(resp.Msg, tt.expectedResp) {
					t.Fatalf("Expected response %v, got %v", tt.expectedResp, resp.Msg)
				}
			}
		})
	}
}

func TestStubService_AllMethods_ReturnNotFound(t *testing.T) {
	service := NewStubService()
	ctx := context.Background()

	// Test all methods except HelloWorld return NOT_FOUND for empty requests
	testCases := []struct {
		name string
		call func() error
	}{
		{
			name: "ProjectNew",
			call: func() error {
				_, err := service.ProjectNew(ctx, connect.NewRequest(&v1.ProjectNewRequest{}))
				return err
			},
		},
		{
			name: "ProjectAssignDiskDirs",
			call: func() error {
				_, err := service.ProjectAssignDiskDirs(ctx, connect.NewRequest(&v1.ProjectAssignDiskDirsRequest{}))
				return err
			},
		},
		{
			name: "ProjectCategorizeFiles",
			call: func() error {
				_, err := service.ProjectCategorizeFiles(ctx, connect.NewRequest(&v1.ProjectCategorizeFilesRequest{}))
				return err
			},
		},
		{
			name: "ProjectSetMetadata",
			call: func() error {
				_, err := service.ProjectSetMetadata(ctx, connect.NewRequest(&v1.ProjectSetMetadataRequest{}))
				return err
			},
		},
		{
			name: "ProjectFinish",
			call: func() error {
				_, err := service.ProjectFinish(ctx, connect.NewRequest(&v1.ProjectFinishRequest{}))
				return err
			},
		},
		{
			name: "ProjectAbandon",
			call: func() error {
				_, err := service.ProjectAbandon(ctx, connect.NewRequest(&v1.ProjectAbandonRequest{}))
				return err
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.call()
			if err == nil {
				t.Fatal("Expected NOT_FOUND error but got none")
			}

			var connectErr *connect.Error
			if !errors.As(err, &connectErr) {
				t.Fatal("Expected connect.Error but got different error type")
			}
			if connectErr.Code() != connect.CodeNotFound {
				t.Fatalf("Expected NOT_FOUND error, got %v", connectErr.Code())
			}
		})
	}
}