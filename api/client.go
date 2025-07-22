package api

import (
	"net/http"
	"time"

	http_client "github.com/lllllan02/http-client"
)

type Client struct {
	*http_client.Client
}

func NewClient(cookie string) (*Client, error) {
	headers := map[string]string{
		"accept-language":    "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6",
		"sec-ch-ua":          `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`,
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "macOS",
		"sec-fetch-site":     "same-origin",
		"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0",
		"referer":            "https://tiku.scratchor.com",
	}

	handler := func(*http.Response) error {
		return nil
	}

	retryStrategy := http_client.DefaultRetryStrategy.
		WithMaxRetries(3).
		WithInterval(5 * time.Second).
		WithRetryCondition(func(resp *http.Response, err error) bool {
			//  如果返回 403，则重试
			if err != nil || resp.StatusCode == http.StatusForbidden {
				return true
			}

			return false
		})

	client := http_client.New(
		http_client.WithCookie("https://tiku.scratchor.com", cookie),
		http_client.WithHeaders(headers),
		http_client.WithLimiter(5),
		http_client.WithResponseHandler(handler),
		http_client.WithRetryStrategy(retryStrategy),
	)

	return &Client{Client: client}, nil
}

func (client *Client) Get(url string) (string, error) {
	return client.Client.Get(url, http_client.WithRequestHeaders(
		map[string]string{
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"priority":                  "u=0, i",
			"sec-fetch-dest":            "document",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-user":            "?1",
			"upgrade-insecure-requests": "1",
		},
	))
}

func (client *Client) Post(url string, data any) (string, error) {
	return client.Client.Post(url, "application/json", data, http_client.WithRequestHeaders(
		map[string]string{
			"origin":           "https://tiku.scratchor.com",
			"priority":         "u=1, i",
			"sec-fetch-dest":   "empty",
			"sec-fetch-mode":   "cors",
			"x-requested-with": "XMLHttpRequest",
		},
	))
}
