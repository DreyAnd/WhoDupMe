package httpclient

import (
	"io"
	"net/http"
	"time"
)

type HttpClient struct {
	client  *http.Client
	headers map[string]string
	cookies []*http.Cookie
}

func NewHttpClient(timeout time.Duration) *HttpClient {
	return &HttpClient{
		client: &http.Client{
			Timeout: timeout,
		},
		headers: make(map[string]string),
	}
}

func (c *HttpClient) SetHeader(key, value string) {
	c.headers[key] = value
}

func (c *HttpClient) SetCookies(cookies []*http.Cookie) {
	c.cookies = cookies
}

func (c *HttpClient) Get(url string) (*http.Response, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// Add custom headers to the request
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Add cookies to the request
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	return c.client.Do(req)
}

func (c *HttpClient) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}

	// Add custom headers to the request
	for key, value := range c.headers {
		req.Header.Set(key, value)
	}

	// Add cookies to the request
	for _, cookie := range c.cookies {
		req.AddCookie(cookie)
	}

	return c.client.Do(req)
}
