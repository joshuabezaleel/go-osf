package osf

import (
	"context"
	"fmt"
)

type FilesService service

type File struct {
	ID string `json:"id"`

	Kind                  string `json:"kind"`
	Name                  string `json:"name"`
	LastTouched           *Time  `json:"last_touched"`
	MaterializedPath      string `json:"materialized_path"`
	DateModified          string `json:"date_modified"`
	CurrentVersion        int64  `json:"current_version"`
	DeleteAllowed         bool   `json:"delete_allowed"`
	DateCreated           *Time  `json:"date_created"`
	Provider              string `json:"provider"`
	Path                  string `json:"path"`
	CurrentUserCanComment bool   `json:"current_user_can_comment"`
	GUID                  string `json:"guid"`
	// Checkout Checkout `json:"checkout"`
	// Tags [][]Tags `json:"tags"`
	Size int64 `json:"size"`

	FileLinks *FileLinks `json:"-"`
}

// FileLinks contains properties inside the link struct that is related to file endpoints according to the Waterbutler API convention.
type FileLinks struct {
	NewFolder *string `json:"new_folder"`
	Move      *string `json:"move"`
	Upload    *string `json:"upload"`
	Download  *string `json:"download"`
	Delete    *string `json:"delete"`
}

func (s *FilesService) GetFileByID(ctx context.Context, id string) (*File, *SingleResponse[*File], error) {
	u := fmt.Sprintf("files/%s", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := doSingle[*File](s.client, ctx, req)
	if err != nil {
		return nil, nil, err
	}

	file := res.GetData()

	// TODO: Do this automatically using reflection.
	if res.Data.Links != nil {
		links := new(FileLinks)
		if i, ok := res.Data.Links["new_folder"]; ok {
			if s, ok := i.(string); ok {
				links.NewFolder = &s
			}
		}
		if i, ok := res.Data.Links["move"]; ok {
			if s, ok := i.(string); ok {
				links.Move = &s
			}
		}
		if i, ok := res.Data.Links["upload"]; ok {
			if s, ok := i.(string); ok {
				links.Upload = &s
			}
		}
		if i, ok := res.Data.Links["download"]; ok {
			if s, ok := i.(string); ok {
				links.Download = &s
			}
		}
		if i, ok := res.Data.Links["delete"]; ok {
			if s, ok := i.(string); ok {
				links.Delete = &s
			}
		}

		file.FileLinks = links
	}

	return file, res, nil
}
