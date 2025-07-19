package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	v1 "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"buf.build/gen/go/krelinga/proto/connectrpc/go/krelinga/video/in/v1/inv1connect"
	"connectrpc.com/connect"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// StubService implements the ServiceHandler interface with empty responses
type StubService struct{}

// HelloWorld returns an empty HelloWorldResponse
func (s *StubService) HelloWorld(ctx context.Context, req *connect.Request[v1.HelloWorldRequest]) (*connect.Response[v1.HelloWorldResponse], error) {
	return connect.NewResponse(&v1.HelloWorldResponse{}), nil
}

// ProjectList returns an empty ProjectListResponse
func (s *StubService) ProjectList(ctx context.Context, req *connect.Request[v1.ProjectListRequest]) (*connect.Response[v1.ProjectListResponse], error) {
	return connect.NewResponse(&v1.ProjectListResponse{}), nil
}

// ProjectNew returns an empty ProjectNewResponse
func (s *StubService) ProjectNew(ctx context.Context, req *connect.Request[v1.ProjectNewRequest]) (*connect.Response[v1.ProjectNewResponse], error) {
	return connect.NewResponse(&v1.ProjectNewResponse{}), nil
}

// UnclaimedDiscDirList returns an empty UnclaimedDiscDirListResponse
func (s *StubService) UnclaimedDiscDirList(ctx context.Context, req *connect.Request[v1.UnclaimedDiscDirListRequest]) (*connect.Response[v1.UnclaimedDiscDirListResponse], error) {
	return connect.NewResponse(&v1.UnclaimedDiscDirListResponse{}), nil
}

// ProjectAssignDiskDirs returns an empty ProjectAssignDiskDirsResponse
func (s *StubService) ProjectAssignDiskDirs(ctx context.Context, req *connect.Request[v1.ProjectAssignDiskDirsRequest]) (*connect.Response[v1.ProjectAssignDiskDirsResponse], error) {
	return connect.NewResponse(&v1.ProjectAssignDiskDirsResponse{}), nil
}

// ProjectGet returns an empty ProjectGetResponse
func (s *StubService) ProjectGet(ctx context.Context, req *connect.Request[v1.ProjectGetRequest]) (*connect.Response[v1.ProjectGetResponse], error) {
	return connect.NewResponse(&v1.ProjectGetResponse{}), nil
}

// ProjectCategorizeFiles returns an empty ProjectCategorizeFilesResponse
func (s *StubService) ProjectCategorizeFiles(ctx context.Context, req *connect.Request[v1.ProjectCategorizeFilesRequest]) (*connect.Response[v1.ProjectCategorizeFilesResponse], error) {
	return connect.NewResponse(&v1.ProjectCategorizeFilesResponse{}), nil
}

// MovieSearch returns an empty MovieSearchResponse
func (s *StubService) MovieSearch(ctx context.Context, req *connect.Request[v1.MovieSearchRequest]) (*connect.Response[v1.MovieSearchResponse], error) {
	return connect.NewResponse(&v1.MovieSearchResponse{}), nil
}

// ProjectSetMetadata returns an empty ProjectSetMetadataResponse
func (s *StubService) ProjectSetMetadata(ctx context.Context, req *connect.Request[v1.ProjectSetMetadataRequest]) (*connect.Response[v1.ProjectSetMetadataResponse], error) {
	return connect.NewResponse(&v1.ProjectSetMetadataResponse{}), nil
}

// ProjectFinish returns an empty ProjectFinishResponse
func (s *StubService) ProjectFinish(ctx context.Context, req *connect.Request[v1.ProjectFinishRequest]) (*connect.Response[v1.ProjectFinishResponse], error) {
	return connect.NewResponse(&v1.ProjectFinishResponse{}), nil
}

// ProjectAbandon returns an empty ProjectAbandonResponse
func (s *StubService) ProjectAbandon(ctx context.Context, req *connect.Request[v1.ProjectAbandonRequest]) (*connect.Response[v1.ProjectAbandonResponse], error) {
	return connect.NewResponse(&v1.ProjectAbandonResponse{}), nil
}

func main() {
	stubService := &StubService{}
	path, handler := inv1connect.NewServiceHandler(stubService)
	
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
