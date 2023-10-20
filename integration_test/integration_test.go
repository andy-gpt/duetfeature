package integration_test

import (
	. "github.com/Eun/go-hit"
	"net/http"
	"testing"
)

// We might do some integration tests here. For example if there was an http interface to search in the database, we could do something like this:

// HTTP GET: /search.
func TestHTTPSearch(t *testing.T) {
	Test(t,
		Description("Database search"),
		Get("/search?q=hello"),
		Expect().Status().Equal(http.StatusOK),
		Expect().Body().String().Contains(`{"results":[{`),
	)
}
