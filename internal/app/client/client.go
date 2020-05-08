package client

import (
	"net/http"

	"github.com/YumaUeno123/linebot_go/internal/app/server/linebot"
)

const (
	UrlScheme      = "https"
	MaxCarouselNum = 10
)

type ClientType int

type Client interface {
	GetKind() string
	Fetch(keyword string) (*[]linebot.Response, error)
}

type Api interface {
	Get(url string) (*http.Response, error)
}

type apiClient struct{}

func (c *apiClient) Get(url string) (*http.Response, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

type retry struct {
	*apiClient
}

func (c *retry) Get(url string) (*http.Response, error) {
	resp, err := c.apiClient.Get(url)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return c.apiClient.Get(url)
	}

	return resp, nil
}

func NewRetry() Api {
	return &retry{
		apiClient: &apiClient{},
	}
}
