package main

import (
	"context"
	"log"
	"os"

	"github.com/joshuabezaleel/go-osf/osf"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("OSF_API_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := osf.NewClient(tc)

	// Fetch 100 first pages, batched 10 per request.
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
			log.Printf("URL: %s", *preprint.Links.Html)
			log.Printf("Title: %s", preprint.Title)
			log.Printf("Description: %s", preprint.Description)
		}

		total += res.PaginationMeta.PerPage
		if total >= res.PaginationMeta.Total {
			break
		}

		opts.Page = res.PaginationMeta.Page+1
	}

	log.Printf("Fetched %d preprints", total)
}