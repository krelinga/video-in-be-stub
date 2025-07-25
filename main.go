package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"

	"buf.build/gen/go/krelinga/proto/connectrpc/go/krelinga/video/in/v1/inv1connect"
	v1 "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// LoggingInterceptor implements connect.Interceptor to log all RPC calls
type LoggingInterceptor struct{}

// WrapUnary implements the Interceptor interface for unary RPC calls
func (l *LoggingInterceptor) WrapUnary(next connect.UnaryFunc) connect.UnaryFunc {
	return func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
		// Get the procedure name from the request
		procedure := req.Spec().Procedure

		// Convert request message to JSON for logging
		reqJSON := ""
		if req.Any() != nil {
			if msg, ok := req.Any().(proto.Message); ok {
				if jsonBytes, err := protojson.Marshal(msg); err == nil {
					reqJSON = string(jsonBytes)
				}
			}
		}

		// Call the actual handler
		resp, err := next(ctx, req)

		// Log the RPC call with error handling
		if err != nil {
			log.Printf("RPC Call [%s] - Request: %s - Error: %v", procedure, reqJSON, err)
		} else {
			// Convert response message to JSON for logging (if successful)
			respJSON := ""
			if resp != nil && resp.Any() != nil {
				if msg, ok := resp.Any().(proto.Message); ok {
					if jsonBytes, errJSON := protojson.Marshal(msg); errJSON == nil {
						respJSON = string(jsonBytes)
					}
				}
			}
			log.Printf("RPC Call [%s] - Request: %s - Response: %s", procedure, reqJSON, respJSON)
		}

		return resp, err
	}
}

// WrapStreamingClient implements the Interceptor interface for streaming client calls
func (l *LoggingInterceptor) WrapStreamingClient(next connect.StreamingClientFunc) connect.StreamingClientFunc {
	return next // No streaming clients in this stub service
}

// WrapStreamingHandler implements the Interceptor interface for streaming handler calls
func (l *LoggingInterceptor) WrapStreamingHandler(next connect.StreamingHandlerFunc) connect.StreamingHandlerFunc {
	return next // No streaming handlers in this stub service
}

// RequestResponseMapping represents a mapping between a request and its corresponding response
type RequestResponseMapping[Req, Resp proto.Message] struct {
	Request  Req
	Response Resp
}

// StubService implements the ServiceHandler interface with configurable responses
type StubService struct {
	// Mappings for each RPC method
	helloWorldMappings             []RequestResponseMapping[*v1.HelloWorldRequest, *v1.HelloWorldResponse]
}

// findMatchingResponse searches for a matching request and returns the corresponding response
func findMatchingResponse[Req, Resp proto.Message](req Req, mappings []RequestResponseMapping[Req, Resp]) (Resp, error) {
	for _, mapping := range mappings {
		if proto.Equal(req, mapping.Request) {
			return mapping.Response, nil
		}
	}
	var zero Resp
	return zero, connect.NewError(connect.CodeNotFound, fmt.Errorf("no matching request found"))
}

// HelloWorld searches for a matching request and returns the corresponding response
func (s *StubService) HelloWorld(ctx context.Context, req *connect.Request[v1.HelloWorldRequest]) (*connect.Response[v1.HelloWorldResponse], error) {
	resp, err := findMatchingResponse(req.Msg, s.helloWorldMappings)
	if err != nil {
		return nil, err
	}
	return connect.NewResponse(resp), nil
}

// ProjectList searches for a matching request and returns the corresponding response
func (*StubService) ProjectList(ctx context.Context, req *connect.Request[v1.ProjectListRequest]) (*connect.Response[v1.ProjectListResponse], error) {
	resp := &v1.ProjectListResponse{}
	for _, p := range data.Projects {
		resp.Projects = append(resp.Projects, p.Project)
	}
	return connect.NewResponse(resp), nil
}

// ProjectNew searches for a matching request and returns the corresponding response
func (s *StubService) ProjectNew(ctx context.Context, req *connect.Request[v1.ProjectNewRequest]) (*connect.Response[v1.ProjectNewResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ProjectNew is not implemented"))
}

// UnclaimedDiscDirList searches for a matching request and returns the corresponding response
func (s *StubService) UnclaimedDiscDirList(ctx context.Context, req *connect.Request[v1.UnclaimedDiscDirListRequest]) (*connect.Response[v1.UnclaimedDiscDirListResponse], error) {
	resp := &v1.UnclaimedDiscDirListResponse{}
	resp.Dirs = append(resp.Dirs, data.Unclaimed...)
	return connect.NewResponse(resp), nil
}

// ProjectAssignDiskDirs searches for a matching request and returns the corresponding response
func (s *StubService) ProjectAssignDiskDirs(ctx context.Context, req *connect.Request[v1.ProjectAssignDiskDirsRequest]) (*connect.Response[v1.ProjectAssignDiskDirsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ProjectAssignDiskDirs is not implemented"))
}

// ProjectGet searches for a matching request and returns the corresponding response
func (s *StubService) ProjectGet(ctx context.Context, req *connect.Request[v1.ProjectGetRequest]) (*connect.Response[v1.ProjectGetResponse], error) {
	found := data.FindProject(req.Msg.Project)
	if found == nil {
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("project not found: %s", req.Msg.Project))
	}
	return connect.NewResponse(found), nil
}

// ProjectCategorizeFiles searches for a matching request and returns the corresponding response
func (s *StubService) ProjectCategorizeFiles(ctx context.Context, req *connect.Request[v1.ProjectCategorizeFilesRequest]) (*connect.Response[v1.ProjectCategorizeFilesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ProjectCategorizeFiles is not implemented"))
}

// MovieSearch searches for a matching request and returns the corresponding response
func (s *StubService) MovieSearch(ctx context.Context, req *connect.Request[v1.MovieSearchRequest]) (*connect.Response[v1.MovieSearchResponse], error) {
	resp := &v1.MovieSearchResponse{}
	results := data.FindMetadata(req.Msg.PartialTitle)
	resp.Results = slices.Collect(results)
	return connect.NewResponse(resp), nil
}

// ProjectSetMetadata searches for a matching request and returns the corresponding response
func (s *StubService) ProjectSetMetadata(ctx context.Context, req *connect.Request[v1.ProjectSetMetadataRequest]) (*connect.Response[v1.ProjectSetMetadataResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ProjectSetMetadata is not implemented"))
}

// ProjectFinish searches for a matching request and returns the corresponding response
func (s *StubService) ProjectFinish(ctx context.Context, req *connect.Request[v1.ProjectFinishRequest]) (*connect.Response[v1.ProjectFinishResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ProjectFinish is not implemented"))
}

// ProjectAbandon searches for a matching request and returns the corresponding response
func (s *StubService) ProjectAbandon(ctx context.Context, req *connect.Request[v1.ProjectAbandonRequest]) (*connect.Response[v1.ProjectAbandonResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, fmt.Errorf("ProjectAbandon is not implemented"))
}

// NewStubService creates a new StubService with predefined request/response mappings
func NewStubService() *StubService {
	return &StubService{
		// Example mapping for HelloWorld
		helloWorldMappings: []RequestResponseMapping[*v1.HelloWorldRequest, *v1.HelloWorldResponse]{
			{
				Request:  &v1.HelloWorldRequest{Name: ""},
				Response: &v1.HelloWorldResponse{Message: "Hello, empty!"},
			},
			{
				Request:  &v1.HelloWorldRequest{Name: "test"},
				Response: &v1.HelloWorldResponse{Message: "Hello, test!"},
			},
			{
				Request:  &v1.HelloWorldRequest{Name: "world"},
				Response: &v1.HelloWorldResponse{Message: "Hello, world!"},
			},
		},
	}
}

func main() {
	stubService := NewStubService()

	// Create the logging interceptor
	loggingInterceptor := &LoggingInterceptor{}

	// Create the handler with the logging interceptor
	path, handler := inv1connect.NewServiceHandler(
		stubService,
		connect.WithInterceptors(loggingInterceptor),
	)

	mux := http.NewServeMux()
	mux.Handle(path, handler)

	// Support HTTP/2 without TLS for development
	server := &http.Server{
		Addr:    ":8080",
		Handler: h2c.NewHandler(mux, &http2.Server{}),
	}

	fmt.Println("Starting video-in stub server on :8080")
	log.Fatal(server.ListenAndServe())
}
