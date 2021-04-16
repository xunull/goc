package resty_client

import "github.com/go-resty/resty/v2"

func (c *RestyClient) Delete(url string) (*resty.Response, error) {
	return c.Client.R().Delete(url)
}
