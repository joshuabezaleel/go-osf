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
	fileID := "553e69248c5e4a219919ea54"

	file, err := client.Files.GetFileByID(ctx, fileID)
	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(file)
}
