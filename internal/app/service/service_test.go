package service_test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/v-smirnov/just-ad-it/internal/app/infrastructure"
	"github.com/v-smirnov/just-ad-it/internal/app/service"
)

const (
	successfulResponseBody   = "some data"
	incorrectPath            = "bad"
	numOfSuccessfulResponses = 2
	responseBodyHash         = "1e50210a0202497fb79bc38b6ade6c34"
)

func TestService_DoRequests(t *testing.T) {
	cases := []struct {
		name            string
		paths           []string
		isErrorExpected bool
	}{
		{
			name:  "all requests are successful",
			paths: []string{"/first", "/second"},
		},
		{
			name:  "some requests are failed",
			paths: []string{"/first", "/second", incorrectPath},
		},
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/first", getSuccessfulDummyResponse)
	mux.HandleFunc("/second", getSuccessfulDummyResponse)

	testServer := httptest.NewServer(mux)
	defer testServer.Close()

	client := infrastructure.NewClient()
	requestService := service.NewService(client, 2)

	for _, testCase := range cases {
		results := requestService.DoRequests(getListOfURLs(testServer.URL, testCase.paths))

		if len(results) != numOfSuccessfulResponses {
			t.Errorf("got unexpected result %q", results)
		}

		for _, result := range results {
			if !strings.HasSuffix(result, responseBodyHash) {
				t.Errorf("got unexpected hashed response body %q", result)
			}
		}
	}
}

func getSuccessfulDummyResponse(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	io.WriteString(w, successfulResponseBody)
}

func getListOfURLs(baseURL string, paths []string) []string {
	urls := make([]string, 0)

	for _, path := range paths {
		if path == incorrectPath {
			urls = append(urls, fmt.Sprintf("%s%s", "https://", path))
		} else {
			urls = append(urls, fmt.Sprintf("%s%s", baseURL, path))
		}
	}

	return urls
}
