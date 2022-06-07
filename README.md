# Open Science Framework (OSF) SDK for Golang

[![Go Version](https://img.shields.io/github/go-mod/go-version/joshuabezaleel/go-osf)](https://github.com/joshuabezaleel/go-osf)
[![GoDoc](https://img.shields.io/static/v1?label=godoc&message=reference&color=blue)](https://pkg.go.dev/github.com/joshuabezaleel/go-osf@v0.0.3/osf)
[![tests](https://github.com/joshuabezaleel/go-osf/actions/workflows/tests.yaml/badge.svg)](https://github.com/joshuabezaleel/go-osf/actions/workflows/tests.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/joshuabezaleel/go-osf)](https://goreportcard.com/report/github.com/joshuabezaleel/go-osf)

go-osf if a Go client library for accessing the [Open Science Framework (OSF)](https://osf.io) [API](https://developer.osf.io/).

## Installation

test change

go-osf makes use of the [generics](https://go.dev/doc/tutorial/generics), so it requires Go v1.18.

To use this library within a Go module project:

```
go get github.com/joshuabezaleel/go-osf
```

## Usage

An OSF access token is required to use this library. You can create one from the [OSF settings page](https://osf.io/settings/tokens). You can [create an account](https://osf.io/register) if you haven't had already.

The following code gets the first 100 public [preprints](https://osf.io/preprints/) from the OSF.

```go
package main

import (
	"context"
	"log"

	"github.com/joshuabezaleel/go-osf/osf"
	"golang.org/x/oauth2"
)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: "your token here ..."},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := osf.NewClient(tc)

	opts := &osf.PreprintsListOptions{
		ListOptions: osf.ListOptions{
			Page:    1,
			PerPage: 100,
		},
	}
	preprints, _, err := client.Preprints.ListPreprints(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}

	for _, preprint := range preprints {
		log.Printf("URL: %s", *preprint.Links.Html)
		log.Printf("Title: %s", preprint.Title)
		log.Printf("Description: %s", preprint.Description)
	}

}
```

Head over to the [examples folder](examples) or [pkg.go.dev](https://pkg.go.dev/github.com/joshuabezaleel/go-osf) for more usage examples.

## License

This library is distributed under the MIT license found in the [LICENSE.txt](LICENSE.txt) file.
