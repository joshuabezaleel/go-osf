package osf

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPreprintsService_List(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	mux.HandleFunc("preprints", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		testFormValues(t, r, values{
			"page":          "2",
			"per_page":      "2",
			"reviews_state": "pending",
		})
		fmt.Fprint(w, `[{"number":1}]`)
	})

	opts := &PreprintsListOptions{
		ListOptions: ListOptions{
			Page:    2,
			PerPage: 2,
			Filter: map[string]string{
				"reviews_state": "pending",
			},
		},
	}
	ctx := context.Background()
	preprints, _, err := client.Preprints.ListPreprints(ctx, opts)
	if err != nil {
		t.Errorf("Preprints.List returned error: %v", err)
	}

	assert.NotZero(t, len(preprints))
	for _, preprint := range preprints {
		assert.Equal(t, "pending", preprint.ReviewsState)
	}

}
