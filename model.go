package main

import (
	"iter"
	"slices"
	"strings"

	v1 "buf.build/gen/go/krelinga/proto/protocolbuffers/go/krelinga/video/in/v1"
	"github.com/krelinga/go-iters"
)

type Model struct {
	Projects  []*v1.ProjectGetResponse
	Unclaimed []string
	Metadata  []*v1.MovieSearchResult
}

func (m *Model) FindProject(name string) *v1.ProjectGetResponse {
	for _, p := range m.Projects {
		if p.Project == name {
			return p
		}
	}
	return nil
}

func (m *Model) FindMetadata(name string) iter.Seq[*v1.MovieSearchResult] {
	name = strings.ToLower(name)
	return iters.Filter(slices.Values(m.Metadata), func(item *v1.MovieSearchResult) bool {
		titleMatches := strings.Contains(strings.ToLower(item.Title), name)
		originalTitleMatches := strings.Contains(item.OriginalTitle, name)
		return titleMatches || originalTitleMatches
	})
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
	Unclaimed: []string{"Unclaimed1", "Unclaimed 2"},
	Metadata: []*v1.MovieSearchResult{
		{
			Title:         "Movie 1",
			OriginalTitle: "Original Movie 1",
			ReleaseDate:   "2023-01-01",
			Genres:        []string{"Action", "Adventure"},
			Overview:      "An action-packed adventure movie.",
		},
		{
			Title:         "Movie 2",
			OriginalTitle: "Original Movie 2",
			ReleaseDate:   "2023-01-02",
			Genres:        []string{"Drama"},
			Overview:      "A dramatic story of love and loss.",
		},
	},
}
