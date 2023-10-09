package ghttp

import (
	"bytes"
	"context"
	"net/http"
	netUrl "net/url"
	"strings"
)

type requestBuilder interface {
	jsonBuild(ctx context.Context, method, url string, request interface{}) (*http.Request, error)
	encodedBuild(ctx context.Context, method, url string, params map[string]string) (*http.Request, error)
}

type httpRequestBuilder struct {
	marshaller marshaller
}

func newRequestBuilder() *httpRequestBuilder {
	return &httpRequestBuilder{
		marshaller: &jsonMarshaller{},
	}
}

func (b *httpRequestBuilder) jsonBuild(ctx context.Context, method, url string, request interface{}) (*http.Request, error) {
	if request == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}

	var reqBytes []byte
	reqBytes, err := b.marshaller.marshal(request)
	if err != nil {
		return nil, err
	}

	return http.NewRequestWithContext(
		ctx,
		method,
		url,
		bytes.NewBuffer(reqBytes),
	)
}

func (b *httpRequestBuilder) encodedBuild(ctx context.Context, method, url string, params map[string]string) (*http.Request, error) {
	if params == nil {
		return http.NewRequestWithContext(ctx, method, url, nil)
	}

	urlValues := netUrl.Values{}
	for key, val := range params {
		urlValues.Add(key, val)
	}

	return http.NewRequestWithContext(
		ctx,
		method,
		url,
		strings.NewReader(urlValues.Encode()),
	)
}
