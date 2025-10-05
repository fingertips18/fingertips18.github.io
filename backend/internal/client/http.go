package client

import "net/http"

type HttpAPI interface {
	Do(req *http.Request) (*http.Response, error)
}

type httpAPI struct {
	client *http.Client
}

func NewHTTPAPI() HttpAPI {
	return &httpAPI{
		client: &http.Client{},
	}
}

func (h *httpAPI) Do(req *http.Request) (*http.Response, error) {
	return h.client.Do(req)
}
