package infrastructure_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/v-smirnov/just-ad-it/internal/app/infrastructure"
)

const (
	successfulURL = "/ok"
	failedURL     = "/fail"

	successfulResponseBody = "some data"
)

func TestClient_Send(t *testing.T) {
	cases := []struct {
		name            string
		isErrorExpected bool
	}{
		{
			name:            "successful request",
			isErrorExpected: false,
		},
		{
			name:            "failed scenario, server error",
			isErrorExpected: true,
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc(successfulURL, getSuccessfulDummyResponse)
	mux.HandleFunc(failedURL, getFailedDummyResponse)

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	client := infrastructure.NewClient()

	for _, testCase := range cases {
		respBody, err := client.Send(resolveURL(testServer.URL, testCase.isErrorExpected), http.MethodGet)

		if testCase.isErrorExpected {
			if string(respBody) != "" {
				t.Errorf("expected empty response body for failed case")
			}
		} else {
			if err != nil {
				t.Errorf("got unexpected error %q", err)
			}

			if string(respBody) != successfulResponseBody {
				t.Errorf("got unexpected response body %q", string(respBody))
			}
		}
	}
}

func getSuccessfulDummyResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, successfulResponseBody)
}

func getFailedDummyResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusInternalServerError)
}

func resolveURL(url string, isErrorExpected bool) string {
	if isErrorExpected {
		return fmt.Sprintf("%s%s", url, failedURL)
	}

	return fmt.Sprintf("%s%s", url, successfulURL)
}
