package ghttp

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Client struct {
	config         ClientConfig
	requestBuilder requestBuilder
}

func NewClient(opts ...Option) *Client {
	config := DefaultConfig()
	client := NewClientWithConfig(config)

	for _, opt := range opts {
		opt.apply(client)
	}

	return client
}

func NewClientWithConfig(config ClientConfig) *Client {
	return &Client{
		config:         config,
		requestBuilder: newRequestBuilder(),
	}
}

func (c *Client) SendRequest(req *http.Request, v interface{}) error {
	res, err := c.config.Client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	byteBody, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != http.StatusOK {
		return c.httpCodeError(req.URL.String(), res.StatusCode, string(byteBody))
	}

	if v != nil {
		if err = json.Unmarshal(byteBody, &v); err != nil {
			return err
		}
	}

	return nil
}

func (c *Client) NewRequest(ctx context.Context, method string, url string, body interface{}) (*http.Request, error) {
	req, err := c.requestBuilder.jsonBuild(ctx, method, c.fullURL(url), body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/json; charset=utf-8")
	}

	return req, nil
}

func (c *Client) NewEncodedRequest(ctx context.Context, method string, url string, params map[string]string) (*http.Request, error) {
	req, err := c.requestBuilder.encodedBuild(ctx, method, c.fullURL(url), params)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json; charset=utf-8")
	contentType := req.Header.Get("Content-Type")
	if contentType == "" {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	return req, nil
}

func (c *Client) httpCodeError(url string, statusCode int, errMsg string) error {
	return fmt.Errorf("http error, url=%v, statusCode=%v, err=%v", url, statusCode, errMsg)
}

func (c *Client) fullURL(suffix string) string {
	return fmt.Sprintf("%s%s", c.config.BaseURL, suffix)
}