package resty_client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-resty/resty/v2"
	"github.com/rs/zerolog/log"
)

type RestyClient struct {
	Client       *resty.Client
	urlParams    map[string]string
	hasUrlParams bool
}

func (c *RestyClient) Init() {
	c.urlParams = make(map[string]string)
}

func NewClient(ops ...RestyOption) *RestyClient {
	client := &RestyClient{
		Client: resty.New(),
	}
	client.Init()
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
	c.Init()
	c.Client.SetHeader("Authorization", "Bear "+token)
	return c
}

func (c *RestyClient) AddHeader(k, v string) {
	c.Client.Header.Add(k, v)
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

func (c *RestyClient) SetUrlParam(key, value string) {
	c.urlParams[key] = value
	c.hasUrlParams = true
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
		fmt.Printf("%v\n", resp)
		fmt.Printf("%s\n",resp.Request.URL)
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
