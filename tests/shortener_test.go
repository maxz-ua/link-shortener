package tests

import (
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/gavv/httpexpect/v2"
	"github.com/stretchr/testify/require"

	"link-shortener/internal/http-server/handlers/url/save"
	"link-shortener/internal/lib/api"
	"link-shortener/internal/lib/random"
)

const (
	host = "localhost:8087"
)

func Test_HappyPath(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	res := e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: random.NewRandomString(10),
		}).
		WithBasicAuth("user", "pass").
		Expect().
		Status(200).
		JSON().Object()

	res.ContainsKey("alias")
	res.Value("alias").NotNull().String().NotEmpty()
}

func Test_SaveRedirect(t *testing.T) {
	testCases := []struct {
		name  string
		url   string
		alias string
		error string
	}{
		{
			name:  "Valid URL",
			url:   gofakeit.URL(),
			alias: gofakeit.Word() + gofakeit.Word(),
		},
		{
			name:  "Invalid URL",
			url:   "invalid_url", // invalid URL format
			alias: gofakeit.Word(),
			error: "field 'URL' must be a valid URL", // Check the specific validation error
		},
		{
			name:  "Empty Alias",
			url:   gofakeit.URL(),
			alias: "",
		},
		{
			name:  "Empty URL",
			url:   "", // Empty URL to test validation
			alias: gofakeit.Word(),
			error: "field 'URL' is required", // Expect validation error for empty URL
		},
		{
			name:  "Invalid Alias",
			url:   gofakeit.URL(),
			alias: "!@#$%^", // Invalid alias (special characters)
			error: "invalid alias (special characters not allowed)",
		},
		// Add more edge cases here
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			u := url.URL{
				Scheme: "http",
				Host:   host,
			}

			e := httpexpect.Default(t, u.String())

			// Save
			resp := e.POST("/url").
				WithJSON(save.Request{
					URL:   tc.url,
					Alias: tc.alias,
				}).
				WithBasicAuth("user", "pass").
				Expect().Status(http.StatusOK).
				JSON().Object()

			if tc.error != "" {
				// Check if the error message is returned correctly
				resp.NotContainsKey("alias")
				resp.Value("error").String().IsEqual(tc.error)
				return
			}

			alias := tc.alias

			if tc.alias != "" {
				// Check that the alias is returned as expected
				resp.Value("alias").String().IsEqual(tc.alias)
			} else {
				// If no alias provided, ensure one is generated
				resp.Value("alias").String().NotEmpty()
				alias = resp.Value("alias").String().Raw()
			}

			// Redirect test
			testRedirect(t, alias, tc.url)

		})
	}
}

func TestCreateAndDeleteURL(t *testing.T) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
	}
	e := httpexpect.Default(t, u.String())

	// Step 1: Create a URL (POST request)
	resp := e.POST("/url").
		WithJSON(save.Request{
			URL:   gofakeit.URL(),
			Alias: gofakeit.Word() + gofakeit.Word(),
		}).
		WithBasicAuth("user", "pass").
		Expect().Status(http.StatusOK).
		JSON().Object()

	// Assert that the URL creation response contains the expected fields
	createdID := int64(resp.Value("id").Number().Raw())
	// Step 2: Perform Delete (DELETE request using the created ID)
	deleteResp := e.DELETE(fmt.Sprintf("/url/%d", createdID)).
		WithBasicAuth("user", "pass").
		Expect().Status(http.StatusOK).
		JSON().Object()

	// Assert that the delete operation was successful
	require.Equal(t, deleteResp.Value("status").Raw(), "OK")

}

func testRedirect(t *testing.T, alias string, urlToRedirect string) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   alias,
	}

	redirectedToURL, err := api.GetRedirect(u.String())
	require.NoError(t, err)

	require.Equal(t, urlToRedirect, redirectedToURL)
}
