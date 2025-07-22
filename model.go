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
			Discs: []*v1.ProjectDisc{
				{
					Disc: "Disc Waiting Thumbs",
					ThumbState: "waiting",
				},
				{
					Disc: "Disc Working Thumbs",
					ThumbState: "working",
				},
				{
					Disc: "Disc Error Thumbs",
					ThumbState: "error",
				},
				{
					Disc: "Disc Done Thumbs",
					ThumbState: "done",
					DiscFiles: []*v1.DiscFile{
						{
							File: "file1.mkv",
							Category: "main_title",
							Thumb: "file1.jpg",
							HumanSize: "1.2 GB",
							HumanDuration: "01:30:00",
							NumChapters: 10,
						},
						{
							File: "file2.mkv",
							Category: "extra",
							Thumb: "file2.jpg",
							HumanSize: "500 MB",
							HumanDuration: "00:45:00",
							NumChapters: 5,
						},
						{
							File: "file3.mkv",
							Category: "trash",
							Thumb: "file3.jpg",
							HumanSize: "300 MB",
							HumanDuration: "00:30:00",
							NumChapters: 3,
						},
						{
							File: "file4.mkv",
							Thumb: "file4.jpg",
							HumanSize: "400 MB",
							HumanDuration: "00:40:00",
							NumChapters: 4,
						},
					},
				},
			},
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
