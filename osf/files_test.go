package osf

import (
	"context"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilesService_GetByID(t *testing.T) {
	client, mux, _, teardown := setup()
	defer teardown()

	fileID := "553e69248c5e4a219919ea54"

	mux.HandleFunc("files/"+fileID, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
	})

	ctx := context.Background()
	file, _, err := client.Files.GetFileByID(ctx, fileID)
	if err != nil {
		t.Errorf("Files.GetFileByID returned error: %v", err)
	}

	assert.NotNil(t, file)
	assert.Equal(t, file.ID, fileID)

}
