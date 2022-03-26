package infrastructure

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

const requestTimeout = time.Second * 5

type Client interface {
	Send(url, method string) ([]byte, error)
}

type client struct{}

func NewClient() Client {
	return &client{}
}

func (c *client) Send(url, method string) ([]byte, error) {
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(make([]byte, 0)))

	//log.Println(fmt.Sprintf("sending [%s] request to url: %s", method, url), )
	resp, err := getHttpClient().Do(req)
	if err != nil {
		return make([]byte, 0), err
	}
	defer func(body io.ReadCloser) {
		err := body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return make([]byte, 0), err
	}

	return body, nil
}

func getHttpClient() *http.Client {
	return &http.Client{
		Timeout:   requestTimeout,
		Transport: http.DefaultTransport,
	}
}
