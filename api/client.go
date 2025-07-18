package api

import (
	"github.com/lllllan02/scratchor-crawler/http"
)

type Client struct {
	*http.Client
}

func NewClient(cookie string) (*Client, error) {
	client, err := http.NewClient("https://tiku.scratchor.com")
	if err != nil {
		return nil, err
	}

	client.SetCookie(cookie)

	return &Client{Client: client}, nil
}

// 检查Cookie状态
func (c *Client) CheckCookies() {
	c.Client.CheckCookies()
}
