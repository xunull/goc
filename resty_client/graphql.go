package resty_client

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/go-resty/resty/v2"
	"reflect"
)

type GraphRequest struct {
	q    string
	vars map[string]interface{}
}

func NewGraphRequest(q string) *GraphRequest {
	req := &GraphRequest{
		q: q,
	}
	return req
}

type GraphRequestOption struct {
	TagName string
}

type Option func(o *GraphRequestOption)

func WithTagName(name string) Option {
	return func(o *GraphRequestOption) {
		o.TagName = name
	}
}

func (req *GraphRequest) VarStruct(i interface{}, opts ...Option) {
	gro := GraphRequestOption{TagName: "json"}
	for _, o := range opts {
		o(&gro)
	}

	s := reflect.TypeOf(i).Elem()
	v := reflect.ValueOf(i)

	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	for i := 0; i < s.NumField(); i++ {
		if name, ok := s.Field(i).Tag.Lookup(gro.TagName); ok {
			req.Var(name, v.Field(i).Interface())
		}
	}

}

func (req *GraphRequest) Var(key string, value interface{}) {
	if req.vars == nil {
		req.vars = make(map[string]interface{})
	}
	req.vars[key] = value
}

// ---------------------------------------------------------------------------------------------------------------------

func (c *RestyClient) PostGraph(req *GraphRequest, url string) (*resty.Response, error) {
	var rb bytes.Buffer
	rbo := struct {
		Query     string                 `json:"query"`
		Variables map[string]interface{} `json:"variables"`
	}{
		Query:     req.q,
		Variables: req.vars,
	}
	if err := json.NewEncoder(&rb).Encode(rbo); err != nil {
		return nil, errors.New("encode graph data has error")
	} else {

	}
	return c.PostJson(rb.String()).Post(url)
}
