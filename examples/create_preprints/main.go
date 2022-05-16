/*
For this example, we use the most cited paper on the earth.

Source: https://pubmed.ncbi.nlm.nih.gov/14907713/
Paper: https://www.jbc.org/article/S0021-9258(19)52451-6/pdf

Before running this example, download the paper first:

    wget "https://www.jbc.org/article/S0021-9258(19)52451-6/pdf" -O paper.pdf

*/

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
		Title:                   osf.StringPointer("PROTEIN MEASUREMENT WITH THE FOLIN PHENOL REAGENT"),
		Description:             osf.StringPointer("Since 1922 when Wu proposed the use of the Folin phenol reagent for the measurement of proteins (l), a number of modified analytical procedures ut.ilizing this reagent have been reported for the determination of proteins in serum (2-G), in antigen-antibody precipitates (7-9), and in insulin (10). Although the reagent would seem to be recommended by its great sensitivity and the simplicity of procedure possible with its use, it has not found great favor for general biochemical purposes. In the belief that this reagent, nevertheless, has considerable merit for certain application, but that its peculiarities and limitations need to be understood for its fullest exploitation, it has been studied with regard t.o effects of variations in pH, time of reaction, and concentration of reactants, permissible levels of reagents commonly used in handling proteins, and interfering subst.ances. Procedures are described for measuring protein in solution or after precipitation wit,h acids or other agents, and for the determination of as little as 0.2 y of protein."),
		OriginalPublicationDate: osf.MustParseTime(time.RFC3339, "1951-05-28T00:00:00.000Z"),
		Subjects: &[][]string{
			{
				"584240da54be81056cecaab0", // Life Sciences
				"584240da54be81056cecac22", // Biochemistry, Biophysics, and Structural Biology
				"584240d954be81056ceca961", // Biochemistry
			},
		},
		Tags:           &[]string{"protein measurement", "folin phenol reagent"},
		HasDataLinks:   osf.NotApplicable,
		HasPreregLinks: osf.NotApplicable,
		HasCOI:         osf.BoolPointer(false),
		IsPublished:    osf.BoolPointer(false),
	}

	file, err := os.Open("paper.pdf")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	preprint, _, err := client.Preprints.CreatePreprint(ctx, req, file)

	if err != nil {
		log.Fatal(err)
	}

	spew.Dump(preprint)
}
