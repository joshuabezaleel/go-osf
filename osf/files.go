package osf

import (
	"context"
	"fmt"
)

type FilesService service

type File struct {
	ID string `jsonapi:"primary,preprints"`

	Kind                  string    `jsonapi:"attr,kind"`
	Name                  string    `jsonapi:"attr,name"`
	LastTouched           Timestamp `jsonapi:"attr,last_touched"`
	MaterializedPath      string    `jsonapi:"attr,materialized_path"`
	DateModified          string    `jsonapi:"attr,date_modified"`
	CurrentVersion        int64     `jsonapi:"attr,current_version"`
	DeleteAllowed         bool      `jsonapi:"attr,delete_allowed"`
	DateCreated           Timestamp `jsonapi:"attr,date_created"`
	Provider              string    `jsonapi:"attr,provider"`
	Path                  string    `jsonapi:"attr,path"`
	CurrentUserCanComment bool      `jsonapi:"attr,current_user_can_comment"`
	GUID                  string    `jsonapi:"attr,guid"`
	// Checkout Checkout `jsonapi:"attr,checkout"`
	// Tags [][]Tags `jsonapi:"attr,tags"`
	Size int64 `jsonapi:"attr,size"`
}

func (s *FilesService) GetFileByID(ctx context.Context, id string) (*File, error) {
	u := fmt.Sprintf("files/%s", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, err
	}

	var file File
	_, err = s.client.Do(ctx, req, &file)
	if err != nil {
		return nil, err
	}

	return &file, nil
}
