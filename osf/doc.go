/*
Package osf provides a client for using the Open Science Framework (OSF, https://osf.io) API.
Since this client uses generics in its implementation, Go v1.18+ is required.

This project follows the similar structure of github.com/google/go-github.

Usage:

	import (
		"github.com/joshuabezaleel/go-osf"
		"golang.org/x/oauth2"
	)

	func main() {
		ctx := context.Background()
		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: "... your access token ..."},
		)
		tc := oauth2.NewClient(ctx, ts)

		client := osf.NewClient(tc)

		// ...
	}

The OSF API token can be obtained from https://osf.io/settings/tokens.

Response Schema

OSF API conforms the JSON API spec v1.0 (https://jsonapi.org/format/1.0/). The
response payload will be in this generic form:

	type Data[T any, U any] struct {
		Type string  `json:"type"`
		ID   *string `json:"id,omitempty"`

		Attributes    T             `json:"attributes,omitempty"`
		Links         U             `json:"links,omitempty"`
		Relationships Relationships `json:"relationships,omitempty"`
	}

	type SinglePayload[T any, U any] struct {
		Data   *Data[T, U] `json:"data,omitempty"`
		Errors Errors      `json:"errors,omitempty"`
	}

	type ManyPayload[T any, U any] struct {
		Data            []*Data[T, U]    `json:"data"`
		PaginationLinks *PaginationLinks `json:"links"`
		Errors          Errors           `json:"errors,omitempty"`

		PaginationMeta  *PaginationMeta `json:"-"`
	}

Nearly all methods will be in this form:

	t, res, err := client.Service.Method(ctx, ...)
	// t is of type T
	// res if of type SinglePayload[T,U] or ManyPayload[T,U]

For most of the times, you won't need the res object, as the first returned
param should already incorporate Attributes, Links, and Relationships. An
exception would be for pagination, which will be explained below.

Pagination

For methods returning multiple objects, the `res` will be of type ManyPayload[T,U].
and PaginationMeta field will contain the pagination info.

	total := 0
	opts := &osf.PreprintsListOptions{
		ListOptions: osf.ListOptions{
			Page:    1,
			PerPage: 10,
		},
	}

	for total < 100 {
		preprints, res, err := client.Preprints.ListPreprints(ctx, opts)
		if err != nil {
			log.Fatal(err)
		}

		for _, preprint := range preprints {
			// Do something with preprint
		}

		total += res.PaginationMeta.PerPage
		if total >= res.PaginationMeta.Total {
			break
		}

		opts.Page = res.PaginationMeta.Page+1
	}

Error Handling

You may check the error returned for type *Errors, which contains more verbose information
about the error.

	t, _, err := client.Service.Method(ctx, ...)

	if err != nil {
		if errs, ok := err.(*osf.Errors); ok {
			// Handle errs
		} else {
			// Handle err as usual
		}
	}

*/
package osf
