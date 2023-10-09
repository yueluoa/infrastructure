package ghttp

import "net/http"

type ClientConfig struct {
	Client  *http.Client
	BaseURL string
}

func DefaultConfig() ClientConfig {
	return ClientConfig{
		Client: http.DefaultClient,
	}
}
