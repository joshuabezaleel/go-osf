package osf

import (
	"context"
)

type PreprintsService service

// TODO: Fix this.
type Timestamp string

type PreprintLicenseRecord struct {
	CopyrightHolders []string `json:"copyright_holders"`
	Year             string   `json:"year"`
}

type Subject struct {
	ID   string `json:"id"`
	Text string `json:"text"`
}

type Preprint struct {
	ID                      string                `jsonapi:"primary,preprints"`
	DateCreated             Timestamp             `jsonapi:"attr,date_created"`
	DateModified            Timestamp             `jsonapi:"attr,date_modified"`
	DatePublished           Timestamp             `jsonapi:"attr,date_published"`
	OriginalPublicationDate Timestamp             `jsonapi:"attr,original_publication_date"`
	DOI                     *string               `jsonapi:"attr,doi"`
	Title                   string                `jsonapi:"attr,title"`
	Description             string                `jsonapi:"attr,description"`
	IsPublished             bool                  `jsonapi:"attr,is_published"`
	IsPreprintOrphan        bool                  `jsonapi:"attr,is_preprint_orphan"`
	LicenseRecord           PreprintLicenseRecord `jsonapi:"attr,license_record"`
	Tags                    []string              `jsonapi:"attr,tags"`
	PreprintDOICreated      Timestamp             `jsonapi:"attr,preprint_doi_created"`
	DateWithdrawn           Timestamp             `jsonapi:"attr,date_withdrawn"`
	Public                  bool                  `jsonapi:"attr,public"`
	ReviewsState            string                `jsonapi:"attr,reviews_state"`
	DateLastTransitioned    Timestamp             `jsonapi:"attr,date_last_transitioned"`
	HasCOI                  bool                  `jsonapi:"attr,has_coi"`
	// Subjects                [][]Subject           `jsonapi:"attr,subjects"`
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

func (s *PreprintsService) ListPreprints(ctx context.Context, opts *PreprintsListOptions) ([]*Preprint, *Response, error) {
	u, err := addOptions("preprints", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	var preprints []*Preprint

	res, err := s.client.Do(ctx, req, &preprints)
	if err != nil {
		return nil, res, err
	}

	return preprints, res, nil
}
