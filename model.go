package main

import v1 "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"

type Model struct {
	Projects []*v1.ProjectGetResponse
	Unclaimed []string
	Metadata []*v1.MovieSearchResult
}

var data = &Model{
	Projects: []*v1.ProjectGetResponse{
		{
			Project: "Empty",
		},
		{
			Project: "Name With Spaces",
		},
	},
}