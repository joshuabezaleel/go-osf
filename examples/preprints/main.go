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

	// opts := &osf.PreprintsListOptions{
	// 	ListOptions: osf.ListOptions{
	// 		Page:    1000,
	// 		PerPage: 2,
	// 	},
	// }
	// preprints, res, err := client.Preprints.ListPreprints(ctx, opts)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// spew.Dump(res)
	// spew.Dump(preprints)

	// preprint, _, err := client.Preprints.GetPreprintByID(ctx, preprints[0].ID)
	preprint, _, err := client.Preprints.GetPreprintByID(ctx, "4rwj5")
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(preprint)
}
