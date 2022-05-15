package osf

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var (
	// For HasDataLinks and HasPreregLinks.
	available     = "available"
	no            = "no"
	notApplicable = "not_applicable"

	Available     = &available
	No            = &no
	NotApplicable = &notApplicable
)

type PreprintsService service

type PreprintLicenseRecord struct {
	CopyrightHolders []string `json:"copyright_holders"`
	Year             string   `json:"year"`
}

type Subject struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type PreprintLinks struct {
	Self        *string `json:"self"`
	Html        *string `json:"html"`
	PreprintDOI *string `json:"preprint_doi"`
}

type Preprint struct {
	ID string `json:"id"`

	DateCreated                 *Time                  `json:"date_created"`
	DateModified                *Time                  `json:"date_modified"`
	DatePublished               *Time                  `json:"date_published"`
	OriginalPublicationDate     *Time                  `json:"original_publication_date"`
	DOI                         *string                `json:"doi"`
	Title                       string                 `json:"title"`
	Description                 string                 `json:"description"`
	IsPublished                 bool                   `json:"is_published"`
	IsPreprintOrphan            bool                   `json:"is_preprint_orphan"`
	LicenseRecord               *PreprintLicenseRecord `json:"license_record"`
	Tags                        []string               `json:"tags"`
	PreprintDOICreated          *Time                  `json:"preprint_doi_created"`
	DateWithdrawn               *Time                  `json:"date_withdrawn"`
	Public                      bool                   `json:"public"`
	ReviewsState                string                 `json:"reviews_state"`
	DateLastTransitioned        *Time                  `json:"date_last_transitioned"`
	HasCOI                      bool                   `json:"has_coi"`
	ConflictOfInterestStatement *string                `json:"conflict_of_interest_statement"`
	Subjects                    [][]*Subject           `json:"subjects"`
	HasDataLinks                *string                `json:"has_data_links"`
	WhyNoData                   *string                `json:"why_no_data"`
	DataLinks                   []string               `json:"data_links"`
	HasPreregLinks              *string                `json:"has_prereg_links"`
	WhyNoPrereg                 *string                `json:"why_no_prereg"`
	PreregLinks                 []string               `json:"prereg_links"`
	PreregLinkInfo              *string                `json:"prereg_link_info"`
	CurrentUserPermissions      []string               `json:"current_user_permissions"`

	Links *PreprintLinks `json:"links"`
}

type PreprintRequest struct {
	PreprintProviderID          string                 `json:"-"`
	Title                       *string                `json:"title,omitempty"`
	Description                 *string                `json:"description,omitempty"`
	IsPublished                 *bool                  `json:"is_published,omitempty"`
	Subjects                    *[][]string            `json:"subjects,omitempty"`
	OriginalPublicationDate     *Time                  `json:"original_publication_date,omitempty"`
	DOI                         *string                `json:"doi,omitempty"`
	LicenseRecord               *PreprintLicenseRecord `json:"license_record,omitempty"`
	Tags                        *[]string              `json:"tags,omitempty"`
	PreprintDOICreated          *Time                  `json:"preprint_doi_created,omitempty"`
	Public                      *bool                  `json:"public,omitempty"`
	HasCOI                      *bool                  `json:"has_coi,omitempty"`
	ConflictOfInterestStatement *string                `json:"conflict_of_interest_statement,omitempty"`
	HasDataLinks                *string                `json:"has_data_links,omitempty"`
	WhyNoData                   *string                `json:"why_no_data,omitempty"`
	DataLinks                   *[]string              `json:"data_links,omitempty"`
	HasPreregLinks              *string                `json:"has_prereg_links,omitempty"`
	WhyNoPrereg                 *string                `json:"why_no_prereg,omitempty"`
	PreregLinks                 *[]string              `json:"prereg_links,omitempty"`
}

type PreprintsListOptions struct {
	ListOptions
}

func buildPreprint(raw *Data[*Preprint, *PreprintLinks]) (*Preprint, error) {
	obj := raw.Attributes
	obj.Links = raw.Links
	return obj, nil
}

func (s *PreprintsService) ListPreprints(ctx context.Context, opts *PreprintsListOptions) ([]*Preprint, *ManyPayload[*Preprint, *PreprintLinks], error) {
	u, err := addOptionsWithFilter("preprints", opts, opts.Filter)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := doMany(s.client, ctx, req, buildPreprint)
	if err != nil {
		return nil, nil, err
	}

	return res.Data, res, nil
}

func (s *PreprintsService) GetPreprintByID(ctx context.Context, id string) (*Preprint, *SinglePayload[*Preprint, *PreprintLinks], error) {
	u := fmt.Sprintf("preprints/%s", id)

	req, err := s.client.NewRequest(http.MethodGet, u, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := doSingle(s.client, ctx, req, buildPreprint)
	if err != nil {
		return nil, nil, err
	}

	return res.Data, res, nil
}

func (s *PreprintsService) CreatePreprint(ctx context.Context, input *PreprintRequest, primaryFile *os.File) (*Preprint, *SinglePayload[*Preprint, *PreprintLinks], error) {
	isPublished := false

	// We cannot set IsPublished to true if primary_files is not uploaded yet,
	// so we should create the preprint with this option unset first.
	if input.IsPublished != nil {
		isPublished = *input.IsPublished
		input.IsPublished = nil
	}

	requestBody := &SinglePayload[*PreprintRequest, interface{}]{
		RawData: &Data[*PreprintRequest, interface{}]{
			Type:       TypePreprints,
			Attributes: input,
			Relationships: Relationships{
				"provider": Relationship{
					Data: &Data[interface{}, interface{}]{
						ID:   &input.PreprintProviderID,
						Type: TypeProviders,
					},
				},
			},
		},
	}

	req, err := s.client.NewRequest(http.MethodPost, "preprints/", requestBody)
	if err != nil {
		return nil, nil, err
	}

	res, err := doSingle(s.client, ctx, req, buildPreprint)
	if err != nil {
		return nil, nil, err
	}

	preprint := res.Data

	// Upload primary file.
	// TODO: Get file upload url from /preprints/:id/files/
	fileUploadURL, err := url.Parse(fmt.Sprintf("https://files.osf.io/v1/resources/%s/providers/osfstorage/", res.Data.ID))
	if err != nil {
		return res.Data, res, err
	}
	values := fileUploadURL.Query()
	values.Add("kind", "file")
	values.Add("name", primaryFile.Name())
	fileUploadURL.RawQuery = values.Encode()

	fileReq, err := http.NewRequest(http.MethodPut, fileUploadURL.String(), primaryFile)
	if err != nil {
		return nil, nil, err
	}
	fileRes, err := doSingle(s.client, ctx, fileReq, buildFile)
	if err != nil {
		return nil, nil, err
	}
	fileID := fileRes.Data.ID
	if strings.HasPrefix(fileID, "osfstorage/") {
		fileID = strings.TrimPrefix(fileID, "osfstorage/")
	}

	// Update preprint to change primary files.
	updatePrimaryFileBody := &SinglePayload[*PreprintRequest, interface{}]{
		RawData: &Data[*PreprintRequest, interface{}]{
			Type: TypePreprints,
			ID:   &preprint.ID,
			Relationships: Relationships{
				"primary_file": Relationship{
					Data: &Data[interface{}, interface{}]{
						ID:   &fileID,
						Type: TypeFiles,
					},
				},
			},
		},
	}

	updatePrimaryFileReq, err := s.client.NewRequest(http.MethodPatch, fmt.Sprintf("preprints/%s/", preprint.ID), updatePrimaryFileBody)
	if err != nil {
		return nil, nil, err
	}

	res, err = doSingle(s.client, ctx, updatePrimaryFileReq, buildPreprint)
	if err != nil {
		return nil, nil, err
	}
	preprint = res.Data

	if isPublished {
		return s.UpdatePreprint(ctx, preprint.ID, &PreprintRequest{IsPublished: &isPublished}, nil)
	}

	return preprint, res, nil
}

func (s *PreprintsService) UpdatePreprint(ctx context.Context, id string, input *PreprintRequest, relationships Relationships) (*Preprint, *SinglePayload[*Preprint, *PreprintLinks], error) {
	body := &SinglePayload[*PreprintRequest, interface{}]{
		RawData: &Data[*PreprintRequest, interface{}]{
			Type:          TypePreprints,
			ID:            &id,
			Attributes:    input,
			Relationships: relationships,
		},
	}

	publishReq, err := s.client.NewRequest(http.MethodPatch, fmt.Sprintf("preprints/%s/", id), body)
	if err != nil {
		return nil, nil, err
	}

	res, err := doSingle(s.client, ctx, publishReq, buildPreprint)
	if err != nil {
		return nil, nil, err
	}

	return res.Data, res, nil
}
