package osf

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
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

func buildFile(raw *Data[*File, *FileLinks]) (*File, error) {
	obj := raw.Attributes
	obj.FileLinks = raw.Links
	return obj, nil
}

func (s *FilesService) GetFileByID(ctx context.Context, id string) (*File, *SinglePayload[*File, *FileLinks], error) {
	u := fmt.Sprintf("files/%s", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := doSingle(s.client, ctx, req, buildFile)
	if err != nil {
		return nil, nil, err
	}

	return res.Data, res, nil
}

func (s *FilesService) DownloadFile(ctx context.Context, dir string, filename string, file *File) error {
	if dir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return err
		}
		dir = wd
	}

	if filename == "" {
		filename = file.Name
	}
	filepath := dir + "/" + filename

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	req, err := http.NewRequest("GET", *file.FileLinks.Download, nil)
	if err != nil {
		return err
	}

	res, err := s.client.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	_, err = io.Copy(out, res.Body)
	if err != nil {
		return err
	}

	return nil
}
