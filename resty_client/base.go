package resty_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type RestyClient struct {
	Client *resty.Client
}

func NewClient(ops ...RestyOption) *RestyClient {
	client := &RestyClient{
		Client: resty.New(),
	}
	op := getDefaultOption(ops...)

	if op.AddTokenHeader {
		if op.TokenLower {
			client.SetLowerToken(op.TokenValue)
		} else {
			client.SetToken(op.TokenValue)
		}
	}

	if op.HostUrl != "" {
		client.SetHostURL(op.HostUrl)
	}

	if op.Proxy != "" {
		client.SetProxy(op.Proxy)
	}

	return client
}

func NewBearClient(token string) *RestyClient {
	c := &RestyClient{
		Client: resty.New(),
	}
	c.Client.SetHeader("Authorization", "Bear "+token)
	return c
}

func (c *RestyClient) SetProxy(proxy string) {
	c.Client.SetProxy(proxy)
}

func (c *RestyClient) RemoveProxy() {
	c.Client.RemoveProxy()
}

func (c *RestyClient) SetHostURL(url string) {
	c.Client.SetHostURL(url)
}

func (c *RestyClient) SetBearer(token string) {
	c.Client.SetHeader("Authorization", fmt.Sprintf("Bearer %s", token))
}

func (c *RestyClient) SetLowerToken(token string) {
	c.Client.SetHeader("Authorization", fmt.Sprintf("token %s", token))
}

func (c *RestyClient) SetToken(token string) {
	c.Client.SetHeader("Authorization", fmt.Sprintf("Token %s", token))
}

func (c *RestyClient) PostJson(data interface{}) *resty.Request {
	return c.Client.R().SetHeader("Content-Type", "application/json").SetBody(data)
}

func (c *RestyClient) PostFormData(data map[string]string) *resty.Request {
	return c.Client.R().SetFormData(data)
}

func (c *RestyClient) PostData(data interface{}) *resty.Request {
	return c.Client.R().SetBody(data)
}

func (c *RestyClient) Ok(resp *resty.Response) bool {
	if resp.IsError() {
		log.Error().Msg(resp.String())
	}
	return resp.IsSuccess()
}

func (c *RestyClient) Json(resp *resty.Response, target interface{}) (interface{}, error) {
	err := json.Unmarshal(resp.Body(), target)
	if err != nil {
		log.Error().Err(err)
	}
	return target, err
}

func (c *RestyClient) CheckUrlExist(url string) (bool, error) {
	if resp, err := c.Client.R().Get(url); err == nil {
		return resp.IsSuccess(), nil
	} else {
		return false, err
	}
}

func (c *RestyClient) FormatJson(resp *resty.Response) (string, error) {
	var str bytes.Buffer
	if err := json.Indent(&str, resp.Body(), "", "    "); err == nil {
		return str.String(), nil
	} else {
		return "", err
	}
}
