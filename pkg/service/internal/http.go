package internal

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

type HttpClient struct {
	Client *http.Client
	Opts   []HttpOption
}

func NewHttpClient() *HttpClient {
	return &HttpClient{
		Client: &http.Client{},
	}
}

func (c *HttpClient) Do(req *Request, v interface{}) (*Response, error) {
	result, err := c.attempt(req)
	if err != nil {
		return nil, err
	}

	resp, err := c.handleResult(result)
	if err != nil {
		return nil, err
	}

	if v != nil {
		if err := json.Unmarshal(resp.Body, v); err != nil {
			return nil, fmt.Errorf("parse response failed: %v", err)
		}
	}
	return resp, nil
}

func (c *HttpClient) attempt(req *Request) (*attemptResult, error) {
	hr, err := req.buildHttpRequest(c.Opts)
	if err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(hr)
	result := &attemptResult{}
	if err != nil {
		result.Err = err
	} else {
		ir, err := newResponse(resp)
		result.Resp = ir
		result.Err = err
	}
	return result, nil
}

func (c *HttpClient) handleResult(result *attemptResult) (*Response, error) {
	if result.Err != nil {
		return nil, fmt.Errorf("make http call failed: %v", result.Err)
	}

	if !c.hasSuccessStatus(result.Resp) {
		return nil, c.newError(result.Resp)
	}

	return result.Resp, nil
}

func (c *HttpClient) hasSuccessStatus(r *Response) bool {
	return r.Status >= http.StatusOK && r.Status < http.StatusNotModified
}

func (c *HttpClient) newError(r *Response) error {
	var respErr struct {
		Error string
	}
	if err := json.Unmarshal(r.Body, &respErr); err != nil {
		return err
	}
	return errors.New(respErr.Error)
}

type attemptResult struct {
	Resp *Response
	Err  error
}

type Request struct {
	Method string
	URL    string
	Body   HttpEntity
	Opts   []HttpOption
}

func (r *Request) buildHttpRequest(opts []HttpOption) (*http.Request, error) {
	var body io.Reader
	if r.Body != nil {
		b, err := r.Body.Bytes()
		if err != nil {
			return nil, err
		}
		body = bytes.NewBuffer(b)
		opts = append(opts, WithHeader("Content-Type", r.Body.Mime()))
	}

	req, err := http.NewRequest(r.Method, r.URL, body)
	if err != nil {
		return nil, err
	}

	opts = append(opts, r.Opts...)
	for _, o := range r.Opts {
		o(req)
	}
	return req, nil
}

type Response struct {
	Status int
	Header http.Header
	Body   []byte
}

func newResponse(resp *http.Response) (*Response, error) {
	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return &Response{
		Status: resp.StatusCode,
		Header: resp.Header,
		Body:   b,
	}, nil
}

type HttpEntity interface {
	Bytes() ([]byte, error)
	Mime() string
}

type jsonEntity struct {
	Value interface{}
}

func NewJsonEntity(v interface{}) *jsonEntity {
	return &jsonEntity{Value: v}
}

func (e *jsonEntity) Bytes() ([]byte, error) {
	return json.Marshal(e.Value)
}

func (e *jsonEntity) Mime() string {
	return "application/json"
}

type HttpOption func(*http.Request)

func WithHeader(key, value string) HttpOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}
