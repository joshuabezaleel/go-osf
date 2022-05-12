package osf

import (
	"context"
	"fmt"
	"io"
	"log"
)

type CitationsService service

type Citations struct {
}

func (s *CitationsService) ListCitationsStyles(ctx context.Context) error {
	u := fmt.Sprint("citations/styles")

	req, err := s.client.NewRequest("GET", u, nil)
	if err != nil {
		return err
	}

	resp, err := s.client.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	bodyString := string(bodyBytes)
	log.Println(bodyString)

	return nil
}
