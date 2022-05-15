package osf

import (
	"context"
	"fmt"
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

	DateCreated             *Time                  `json:"date_created"`
	DateModified            *Time                  `json:"date_modified"`
	DatePublished           *Time                  `json:"date_published"`
	OriginalPublicationDate *Time                  `json:"original_publication_date"`
	DOI                     *string                `json:"doi"`
	Title                   string                 `json:"title"`
	Description             string                 `json:"description"`
	IsPublished             bool                   `json:"is_published"`
	IsPreprintOrphan        bool                   `json:"is_preprint_orphan"`
	LicenseRecord           *PreprintLicenseRecord `json:"license_record"`
	Tags                    []string               `json:"tags"`
	PreprintDOICreated      *Time                  `json:"preprint_doi_created"`
	DateWithdrawn           *Time                  `json:"date_withdrawn"`
	Public                  bool                   `json:"public"`
	ReviewsState            string                 `json:"reviews_state"`
	DateLastTransitioned    *Time                  `json:"date_last_transitioned"`
	HasCOI                  bool                   `json:"has_coi"`
	Subjects                [][]*Subject           `json:"subjects"`

	Links *PreprintLinks `json:"links"`
}

/*
Skipped attrs:
"current_user_permissions": [],
"conflict_of_interest_statement": null,
"has_data_links": "not_applicable",
"why_no_data": null,
"data_links": [],
"has_prereg_links": "not_applicable",
"why_no_prereg": null,
"prereg_links": [],
"prereg_link_info": "",
*/

type PreprintsListOptions struct {
	ListOptions
}

func buildPreprint(raw *Data[*Preprint, *PreprintLinks]) (*Preprint, error) {
	obj := raw.Attributes
	obj.Links = raw.Links
	return obj, nil
}

func (s *PreprintsService) ListPreprints(ctx context.Context, opts *PreprintsListOptions) ([]*Preprint, *ManyResponse[*Preprint, *PreprintLinks], error) {
	u, err := addOptions("preprints", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := doMany(s.client, ctx, req, buildPreprint)
	if err != nil {
		return nil, nil, err
	}

	return res.Data, res, nil
}

func (s *PreprintsService) GetPreprintByID(ctx context.Context, id string) (*Preprint, *SingleResponse[*Preprint, *PreprintLinks], error) {
	u := fmt.Sprintf("preprints/%s", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := doSingle(s.client, ctx, req, buildPreprint)
	if err != nil {
		return nil, nil, err
	}

	return res.Data, res, nil
}
