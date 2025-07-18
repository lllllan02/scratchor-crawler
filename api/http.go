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

func Get(url, cookie string) (string, error) {
	// 等待限流器允许请求
	globalLimiter.wait()

	// 创建 HTTP 客户端
	client := &http.Client{}

	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置 headers
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("priority", "u=0, i")
	req.Header.Set("referer", "https://tiku.scratchor.com/question/cat/1/list?chapterId=1")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")

	// 设置 Cookie
	req.Header.Set("Cookie", cookie)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return string(body), nil
}

func Post(url, cookie string, data any) (string, error) {
	// 等待限流器允许请求
	globalLimiter.wait()

	var (
		client   = &http.Client{}
		jsonData []byte
		err      error
	)

	// JSON 序列化数据
	if data != nil {
		if jsonData, err = json.Marshal(data); err != nil {
			return "", fmt.Errorf("JSON 序列化失败: %v", err)
		}
	}

	// 创建请求
	req, err := http.NewRequest("POST", url, strings.NewReader(string(jsonData)))
	if err != nil {
		return "", fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置 headers
	req.Header.Set("accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9,en;q=0.8,en-GB;q=0.7,en-US;q=0.6")
	req.Header.Set("content-type", "application/json; charset=UTF-8")
	req.Header.Set("origin", "https://tiku.scratchor.com")
	req.Header.Set("priority", "u=1, i")
	req.Header.Set("referer", "https://tiku.scratchor.com/question/view/abjyobpc2yhjb2ei")
	req.Header.Set("sec-ch-ua", `"Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", "macOS")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0")
	req.Header.Set("x-requested-with", "XMLHttpRequest")

	// 设置 Cookie
	req.Header.Set("Cookie", cookie)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应内容
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("读取响应失败: %v", err)
	}

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("请求失败，状态码: %d", resp.StatusCode)
	}

	return string(body), nil
}
