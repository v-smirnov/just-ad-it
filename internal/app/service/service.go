package service

import (
	"crypto/md5"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/v-smirnov/just-ad-it/internal/app/infrastructure"
)

type Service interface {
	DoRequests(urlList []string) []string
}

type service struct {
	client              infrastructure.Client
	maxParallelRequests int8
}

func NewService(client infrastructure.Client, maxParallelRequests int8) Service {
	return &service{
		client:              client,
		maxParallelRequests: maxParallelRequests,
	}
}

func (s *service) DoRequests(urlList []string) []string {
	wg := &sync.WaitGroup{}

	// limiter allow us to limit amount of request. As soon as channel buffer is full app will wait until
	// any goroutine will finish it's work
	limiter := make(chan struct{}, s.maxParallelRequests)

	// collect results here
	responses := make(chan string, len(urlList))

	for _, url := range urlList {
		wg.Add(1)
		limiter <- struct{}{}
		go s.sendRequest(wg, limiter, normalizeUrlIfNecessary(url), responses)
	}

	wg.Wait()
	close(responses)

	results := make([]string, 0)
	for response := range responses {
		results = append(results, response)
	}

	return results
}

func (s *service) sendRequest(wg *sync.WaitGroup, limiter chan struct{}, url string, responses chan string) {
	defer func() {
		wg.Done()
		<-limiter
	}()

	respBody, err := s.client.Send(url, http.MethodGet)

	if err != nil {
		fmt.Println(err)
		return
	}

	responses <- fmt.Sprintf("%s %x", url, md5.Sum(respBody))
}

func normalizeUrlIfNecessary(url string) string {
	if strings.HasPrefix(url, "http") {
		return url
	}

	return fmt.Sprintf("https://%s", url)
}
