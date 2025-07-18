package crawler

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// 定义专门的错误类型
var (
	ErrDailyLimitExceeded = errors.New("每日查看题目次数已达上限，请升级会员")
	ErrAPIError           = errors.New("API请求失败")
)

func GetAnswer(id, cookie string) (string, error) {
	url := fmt.Sprintf("https://tiku.scratchor.com/question/answer/%s", id)

	// 获取 html
	body, err := Post(url, cookie, nil)
	if err != nil {
		return "", err
	}

	// 解析 JSON 响应
	html, err := parseAnswerResponse(body)
	if err != nil {
		return "", err
	}

	// 解析 html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	// 从 JSON 响应中提取答案
	answerSelection := doc.Find("div.answer-body")
	if answerSelection.Length() == 0 {
		return "", fmt.Errorf("未找到答案内容")
	}

	// 获取答案文本并去除空白字符
	answer := strings.TrimSpace(answerSelection.Text())
	if answer == "" {
		return "", fmt.Errorf("答案内容为空")
	}

	return answer, nil
}

// parseAnswerResponse 解析答案API的JSON响应
func parseAnswerResponse(body string) (string, error) {
	type Response struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			HTML string `json:"html"`
		} `json:"data"`
		Redirect string `json:"redirect,omitempty"`
	}

	var resp Response
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		return "", fmt.Errorf("解析 JSON 响应失败: %v", err)
	}

	// 检查响应码
	if resp.Code != 0 {
		// 处理特定的错误码
		switch resp.Code {
		case -1:
			return "", ErrDailyLimitExceeded
		default:
			return "", fmt.Errorf("%w (code: %d): %s", ErrAPIError, resp.Code, resp.Msg)
		}
	}

	// 检查是否有HTML数据
	if resp.Data.HTML == "" {
		return "", fmt.Errorf("响应中没有HTML数据")
	}

	return resp.Data.HTML, nil
}
