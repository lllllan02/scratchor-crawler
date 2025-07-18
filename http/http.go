package http

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"
)

// HTTP 客户端，支持 cookie 管理
type Client struct {
	httpClient *http.Client
	jar        *cookiejar.Jar
	baseURL    string
	limiter    *rateLimiter
}

// 创建新的 HTTP 客户端
func NewClient(baseURL string) (*Client, error) {
	// 创建 cookie jar
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, fmt.Errorf("创建 cookie jar 失败: %w", err)
	}

	// 创建 HTTP 客户端，增加超时设置
	httpClient := &http.Client{
		Jar:     jar,
		Timeout: 30 * time.Second,
	}

	return &Client{
		httpClient: httpClient,
		jar:        jar,
		baseURL:    baseURL,
		limiter:    newRateLimiter(1), // 降低到每秒1次请求
	}, nil
}

// 设置初始cookie
func (c *Client) SetCookie(cookie string) error {
	if cookie == "" {
		return nil
	}

	// 解析 baseURL
	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("解析 baseURL 失败: %w", err)
	}

	// 解析 cookie 字符串并设置到 jar 中
	cookies := parseCookieString(cookie)
	c.jar.SetCookies(parsedURL, cookies)

	// 打印调试信息
	fmt.Printf("设置Cookie成功，共设置 %d 个cookie\n", len(cookies))
	for _, cookie := range cookies {
		fmt.Printf("  - %s: %s\n", cookie.Name, cookie.Value)
	}

	return nil
}

// 解析 cookie 字符串为 Cookie 对象
func parseCookieString(cookieStr string) []*http.Cookie {
	var cookies []*http.Cookie

	// 分割 cookie 字符串
	pairs := strings.Split(cookieStr, ";")
	for _, pair := range pairs {
		pair = strings.TrimSpace(pair)
		if pair == "" {
			continue
		}

		// 分割 name=value
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) != 2 {
			continue
		}

		name := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		cookie := &http.Cookie{
			Name:  name,
			Value: value,
		}
		cookies = append(cookies, cookie)
	}

	return cookies
}

// 请求配置结构体
type requestConfig struct {
	method      string
	url         string
	data        any
	headers     map[string]string
	contentType string
}

// 设置通用headers
func (c *Client) setCommonHeaders(req *http.Request, config *requestConfig) {
	// 设置通用headers
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")
	req.Header.Set("referer", c.baseURL)

	// 设置自定义 headers
	for key, value := range config.headers {
		req.Header.Set(key, value)
	}

	// 设置 Content-Type
	if config.contentType != "" {
		req.Header.Set("content-type", config.contentType)
	}
}

// 执行 HTTP 请求的通用方法，带重试机制
func (c *Client) executeRequest(config *requestConfig) (string, error) {
	maxRetries := 3
	var lastErr error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		// 等待限流器允许请求
		c.limiter.wait()

		// 准备请求体
		var body io.Reader
		if config.data != nil {
			if jsonData, err := json.Marshal(config.data); err != nil {
				return "", fmt.Errorf("JSON 序列化失败: %w", err)
			} else {
				body = strings.NewReader(string(jsonData))
			}
		}

		// 创建请求
		req, err := http.NewRequest(config.method, config.url, body)
		if err != nil {
			return "", fmt.Errorf("创建请求失败: %w", err)
		}

		// 设置headers
		c.setCommonHeaders(req, config)

		// 发送请求
		resp, err := c.httpClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("发送请求失败: %w", err)
			fmt.Printf("请求失败 (尝试 %d/%d): %v\n", attempt, maxRetries, err)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * 2 * time.Second) // 递增延迟
				continue
			}
			return "", lastErr
		}
		defer resp.Body.Close()

		// 读取响应内容
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			lastErr = fmt.Errorf("读取响应失败: %w", err)
			fmt.Printf("读取响应失败 (尝试 %d/%d): %v\n", attempt, maxRetries, err)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * 2 * time.Second)
				continue
			}
			return "", lastErr
		}

		// 检查响应状态码
		if resp.StatusCode == http.StatusForbidden {
			lastErr = fmt.Errorf("请求被禁止 (403)，可能是Cookie过期或请求频率过高")
			fmt.Printf("收到403错误 (尝试 %d/%d): %s\n", attempt, maxRetries, string(respBody))

			// 如果是403错误，增加更长的延迟
			if attempt < maxRetries {
				delay := time.Duration(attempt) * 5 * time.Second
				fmt.Printf("等待 %v 后重试...\n", delay)
				time.Sleep(delay)
				continue
			}
			return "", lastErr
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(respBody))
			fmt.Printf("请求失败 (尝试 %d/%d): 状态码 %d\n", attempt, maxRetries, resp.StatusCode)
			if attempt < maxRetries {
				time.Sleep(time.Duration(attempt) * 2 * time.Second)
				continue
			}
			return "", lastErr
		}

		// 请求成功
		if attempt > 1 {
			fmt.Printf("请求成功 (重试 %d 次后)\n", attempt-1)
		}
		return string(respBody), nil
	}

	return "", lastErr
}

// GET请求
func (c *Client) Get(url string) (string, error) {
	config := &requestConfig{
		method: "GET",
		url:    url,
		headers: map[string]string{
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"priority":                  "u=0, i",
			"sec-fetch-dest":            "document",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-user":            "?1",
			"upgrade-insecure-requests": "1",
		},
	}

	return c.executeRequest(config)
}

// POST请求
func (c *Client) Post(url string, data any) (string, error) {
	config := &requestConfig{
		method:      "POST",
		url:         url,
		data:        data,
		contentType: "application/json; charset=UTF-8",
		headers: map[string]string{
			"accept":           "application/json, text/javascript, */*; q=0.01",
			"origin":           "https://tiku.scratchor.com",
			"priority":         "u=1, i",
			"sec-fetch-dest":   "empty",
			"sec-fetch-mode":   "cors",
			"x-requested-with": "XMLHttpRequest",
		},
	}

	return c.executeRequest(config)
}

// 令牌桶限流器
type rateLimiter struct {
	tokens chan struct{}
	rate   time.Duration
}

// 创建新的限流器，每秒允许指定数量的请求
func newRateLimiter(qps int) *rateLimiter {
	rl := &rateLimiter{
		tokens: make(chan struct{}, qps),
		rate:   time.Second / time.Duration(qps),
	}

	// 初始化令牌桶
	for range qps {
		rl.tokens <- struct{}{}
	}

	// 启动令牌补充协程
	go rl.refillTokens()

	return rl
}

// 补充令牌
func (rl *rateLimiter) refillTokens() {
	ticker := time.NewTicker(rl.rate)
	defer ticker.Stop()

	for range ticker.C {
		select {
		case rl.tokens <- struct{}{}:
			// 成功添加令牌
		default:
			// 桶已满，跳过
		}
	}
}

// 等待获取令牌
func (rl *rateLimiter) wait() {
	<-rl.tokens
}

// 检查当前Cookie状态
func (c *Client) CheckCookies() {
	parsedURL, err := url.Parse(c.baseURL)
	if err != nil {
		fmt.Printf("解析URL失败: %v\n", err)
		return
	}

	cookies := c.jar.Cookies(parsedURL)
	fmt.Printf("当前Cookie状态 (共 %d 个):\n", len(cookies))
	for _, cookie := range cookies {
		fmt.Printf("  - %s: %s (过期时间: %v)\n", cookie.Name, cookie.Value, cookie.Expires)
	}
}
