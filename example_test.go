package vectara_test

import (
	"context"
	"net/http"
	"os"

	"github.com/psyb0t/go-vectara"
)

//nolint:testableexamples // Needs live credentials, compile check only.
func Example() {
	auth := vectara.WithRequestEditorFn(
		func(_ context.Context, req *http.Request) error {
			req.Header.Set("x-api-key", os.Getenv("VECTARA_API_KEY"))

			return nil
		},
	)

	client, err := vectara.NewClientWithResponses(
		"https://api.vectara.io",
		auth,
	)
	if err != nil {
		panic(err)
	}

	res, err := client.ListCorporaWithResponse(context.Background(), nil)
	if err != nil {
		panic(err)
	}

	_ = res.JSON200
}
