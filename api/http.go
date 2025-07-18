package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

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

// 全局限流器实例，每秒2次请求
var globalLimiter = newRateLimiter(2)

// 请求配置结构体
type requestConfig struct {
	method      string
	url         string
	cookie      string
	data        any
	headers     map[string]string
	contentType string
}

// 设置通用headers
func setCommonHeaders(req *http.Request, config *requestConfig) {
	// 设置通用headers
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")
	req.Header.Set("referer", config.url)

	// 设置Cookie
	req.Header.Set("Cookie", config.cookie)

	// 设置自定义headers
	for key, value := range config.headers {
		req.Header.Set(key, value)
	}

	// 设置Content-Type
	if config.contentType != "" {
		req.Header.Set("content-type", config.contentType)
	}
}

// 执行HTTP请求的通用方法
func executeRequest(config *requestConfig) (string, error) {
	// 等待限流器允许请求
	globalLimiter.wait()

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 准备请求体
	var body io.Reader
	if config.data != nil {
		if jsonData, err := json.Marshal(config.data); err != nil {
			return "", fmt.Errorf("JSON 序列化失败: %v", err)
		} else {
			body = strings.NewReader(string(jsonData))
		}
	}

	// 创建请求
	req, err := http.NewRequest(config.method, config.url, body)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置headers
	setCommonHeaders(req, config)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return string(respBody), nil
}

func Get(url, cookie string) (string, error) {
	config := &requestConfig{
		method: "GET",
		url:    url,
		cookie: cookie,
		headers: map[string]string{
			"accept":                    "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7",
			"priority":                  "u=0, i",
			"sec-fetch-dest":            "document",
			"sec-fetch-mode":            "navigate",
			"sec-fetch-user":            "?1",
			"upgrade-insecure-requests": "1",
		},
	}

	return executeRequest(config)
}

func Post(url, cookie string, data any) (string, error) {
	config := &requestConfig{
		method:      "POST",
		url:         url,
		cookie:      cookie,
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

	return executeRequest(config)
}
