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

func ParseTime(layout, value string) (*Time, error) {
	t, err := time.Parse(layout, value)
	if err != nil {
		return nil, err
	}
	return &Time{t}, nil
}

func MustParseTime(layout, value string) *Time {
	t, err := ParseTime(layout, value)
	if err != nil {
		panic(err)
	}
	return t
}
