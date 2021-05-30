package resty_client

import "github.com/go-resty/resty/v2"

func (c *RestyClient) PostData(data interface{}) *resty.Request {
	return c.Client.R().SetBody(data)
}

func (c *RestyClient) PostJson(data interface{}) *resty.Request {
	return c.Client.R().SetHeader("Content-Type", "application/json").SetBody(data)
}

func (c *RestyClient) PostFormData(data map[string]string) *resty.Request {
	return c.Client.R().SetFormData(data)
}
