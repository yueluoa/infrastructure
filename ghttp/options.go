package ghttp

type Option interface {
	apply(*Client)
}

type LogOption struct {
	f func(*Client)
}

func (lo *LogOption) apply(c *Client) {
	lo.f(c)
}

func NewLogOption(f func(*Client)) *LogOption {
	return &LogOption{
		f: f,
	}
}

func WithBaseURL(baseURL string) Option {
	return NewLogOption(func(c *Client) {
		c.config.BaseURL = baseURL
	})
}
