package osf

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
)

type PreprintProvidersService service

type PreprintProviderLinks struct {
	Self        *string `json:"self"`
	Preprints   *string `json:"preprints"`
	ExternalURL *string `json:"external_url"`
}

type PreprintProviderSubject struct {
	TaxonomiesID       []string
	IncludeAllChildren bool
}

func (s *PreprintProviderSubject) UnmarshalJSON(b []byte) error {
	var obj []interface{}
	if err := json.Unmarshal(b, &obj); err != nil {
		return err
	}
	if len(obj) != 2 {
		return errors.New("the length of a subjects_acceptable element must be 2")
	}
	subjectsIface, ok := obj[0].([]interface{})
	if !ok {
		return errors.New("the first element of a subjects_acceptable element must be an array of string")
	}
	include, ok := obj[1].(bool)
	if !ok {
		return errors.New("the second element of a subjects_acceptable element must be a boolean")
	}
	subjects := make([]string, len(subjectsIface))
	for i, subjectIface := range subjectsIface {
		subject, ok := subjectIface.(string)
		if !ok {
			return errors.New("the first element of a subjects_acceptable element must be an array of string")
		}
		subjects[i] = subject
	}
	s.TaxonomiesID = subjects
	s.IncludeAllChildren = include
	return nil
}

func (s PreprintProviderSubject) MarshalJSON() ([]byte, error) {
	return json.Marshal([]interface{}{s.TaxonomiesID, s.IncludeAllChildren})
}

type PreprintProvider struct {
	ID string `json:"id"`

	Name                    string                     `json:"name"`
	Description             string                     `json:"description"`
	AdvisoryBoard           *string                    `json:"advisory_board"`
	Example                 *string                    `json:"example"`
	Domain                  *string                    `json:"domain"`
	DomainRedirectEnabled   bool                       `json:"domain_redirect_enabled"`
	FooterLinks             *string                    `json:"footer_links"`
	EmailSupport            *string                    `json:"email_support"`
	FacebookAppID           *int64                     `json:"facebook_app_id"`
	AllowSubmissions        bool                       `json:"allow_submissions"`
	AllowCommenting         bool                       `json:"allow_commenting"`
	Assets                  map[string]interface{}     `json:"assets"`
	ShareSource             *string                    `json:"share_source"`
	SharePublishType        *string                    `json:"share_publish_type"`
	Permissions             []string                   `json:"permissions"`
	PreprintWord            *string                    `json:"preprint_word"`
	AdditionalProviders     []string                   `json:"additional_providers"`
	ReviewsWorkflow         *string                    `json:"reviews_workflow"`
	ReviewsCommentPrivate   bool                       `json:"reviews_comment_private"`
	ReviewsCommentAnonymous bool                       `json:"reviews_comment_anonymous"`
	HeaderText              *string                    `json:"header_text"`
	BannerPath              *string                    `json:"banner_path"`
	LogoPath                *string                    `json:"logo_path"`
	EmailContact            *string                    `json:"email_contact"`
	SocialTwitter           *string                    `json:"social_twitter"`
	SocialFacebook          *string                    `json:"social_facebook"`
	SocialInstagram         *string                    `json:"social-instagram"`
	SubjectsAcceptable      []*PreprintProviderSubject `json:"subjects_acceptable"`

	Links *PreprintProviderLinks `json:"links"`
}

type PreprintProvidersListOptions struct {
	ListOptions
}

func buildPreprintProvider(raw *Data[*PreprintProvider, *PreprintProviderLinks]) (*PreprintProvider, error) {
	obj := raw.Attributes
	obj.Links = raw.Links
	return obj, nil
}

func (s *PreprintProvidersService) ListPreprintProviders(ctx context.Context, opts *PreprintProvidersListOptions) ([]*PreprintProvider, *ManyResponse[*PreprintProvider, *PreprintProviderLinks], error) {
	u, err := addOptions("preprint_providers", opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := doMany(s.client, ctx, req, buildPreprintProvider)
	if err != nil {
		return nil, nil, err
	}

	return res.Data, res, nil
}

func (s *PreprintProvidersService) GetPreprintProviderByID(ctx context.Context, id string) (*PreprintProvider, *SingleResponse[*PreprintProvider, *PreprintProviderLinks], error) {
	u := fmt.Sprintf("preprint_providers/%s", id)

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	res, err := doSingle(s.client, ctx, req, buildPreprintProvider)
	if err != nil {
		return nil, nil, err
	}

	return res.Data, res, nil
}
