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

	file, _, err := client.Preprints.GetPreprintPrimaryFileByID(ctx, "xfdsr")
	if err != nil {
		log.Fatal(err)
	}

	err = client.Files.DownloadFile(ctx, "", "", file)
	if err != nil {
		log.Fatal(err)
	}
}
