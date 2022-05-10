package examples

import (
	"context"
	"log"
	"os"

	goosf "github.com/joshuabezaleel/gohack-osfapi-wrapper/osf"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("OSF_API_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := goosf.NewClient(tc)
	err := client.Citations.ListCitationsStyles(ctx)
	if err != nil {
		log.Fatal(err)
	}

}
