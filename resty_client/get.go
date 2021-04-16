package resty_client

import (
	"fmt"
	"github.com/go-resty/resty/v2"
)

func (c *RestyClient) Get(url string) (*resty.Response, error) {
	if c.hasUrlParams {
		for key, value := range c.urlParams {
			url = fmt.Sprintf("%s?%s=%s", url, key, value)
		}
	}
	return c.Client.R().Get(url)
}

func (c *RestyClient) GetWithQueryMap(url string, query map[string]string) (*resty.Response, error) {
	return c.Client.R().SetQueryParams(query).Get(url)
}
