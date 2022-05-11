package main

import (
	"context"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
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
	opts := &osf.PreprintsListOptions{
		ListOptions: osf.ListOptions{
			Page: 1000,
			PerPage: 2,
		},
	}
	preprints, err := client.Preprints.ListPreprints(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(preprints)
}
