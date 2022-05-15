package main

import (
	"context"
	"log"
	"os"
	"time"

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

	req := &osf.PreprintRequest{
		PreprintProviderID:      "osf",
		Title:                   osf.StringPointer("Minimax and Expectimax Algorithm to Solve 2048"),
		Description:             osf.StringPointer("2048 is a puzzle game created by Gabriele Cirulli a few months ago. It was booming recently and played by millions of people over the internet. People keep searching for the optimal algorithm for solving the game. Here are few approaches: minimax and expectimax algorithm. The idea is to calculate all possible moves and then select the best move by some functions. Alpha-beta pruning is also used to speed up search time. The result depends on the limit of the depth of the search tree. The greater the limit, the better the result. At some point, expectimax algorithm reaches 80% winning rate."),
		OriginalPublicationDate: osf.MustParseTime(time.RFC3339, "2018-05-18T00:00:00.000Z"),
		Subjects: &[][]string{
			{
				"584240d954be81056ceca9a1", // Physical Sciences and Mathematics
				"584240da54be81056cecabbb", // Computer Sciences
				"584240da54be81056cecaa84", // Theory and Algorithms
			},
		},
		Tags: &[]string{
			"2048",
			"Depth Limited Search",
			"Expectimax",
			"Heuristic Function",
			"Minimax",
			"Search Tree",
		},
		HasDataLinks:   osf.Available,
		DataLinks:      &[]string{"https://play2048.co/"},
		HasPreregLinks: osf.NotApplicable,
		HasCOI:         osf.BoolPointer(false),
		IsPublished:    osf.BoolPointer(true),
	}

	file, err := os.Open("paper.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	preprint, res, err := client.Preprints.CreatePreprint(ctx, req, file)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(res)
	spew.Dump(preprint)
}
