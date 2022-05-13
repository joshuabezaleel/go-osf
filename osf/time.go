package osf

import (
	"strings"
	"time"
)

type Time struct {
	time.Time
}

func (t *Time) UnmarshalJSON(b []byte) error {
	// Strip quotes.
	s := string(b)
	s = strings.Trim(s, `"`)

	var parsed time.Time
	var err error

	if strings.Contains(s, "Z") {
		parsed, err = time.Parse(`"2006-01-02T15:04:05Z"`, string(b))
	} else {
		parsed, err = time.Parse(`"2006-01-02T15:04:05"`, string(b))
	}

	if err != nil {
		return err
	}

	*t = Time{parsed}

	return nil
}
