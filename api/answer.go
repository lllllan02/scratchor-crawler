package api

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

// AnswerResult 表示答案和解析的结果
type AnswerResult struct {
	Answer   []string `json:"answer"`
	Analysis string   `json:"analysis"`
}

func (client *Client) GetAnswer(id string, questionType string) (*AnswerResult, error) {
	url := fmt.Sprintf("https://tiku.scratchor.com/question/answer/%s", id)

	// 获取 html
	body, err := client.Post(url, nil)
	if err != nil {
		return nil, err
	}

	// 解析 JSON 响应
	html, err := parseAnswerResponse(body)
	if err != nil {
		return nil, err
	}

	// 解析 html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, err
	}

	// 提取答案
	answer, _ := extractAnswer(doc, questionType)

	// 提取解析
	analysis, _ := extractAnalysis(doc)

	return &AnswerResult{Answer: answer, Analysis: analysis}, nil
}

// extractAnswer 根据题型提取答案内容
func extractAnswer(doc *goquery.Document, questionType string) ([]string, error) {
	answerSelection := doc.Find("div.answer-body")
	if answerSelection.Length() == 0 {
		return nil, fmt.Errorf("未找到答案内容")
	}

	// 根据题型进行不同的解析逻辑
	switch questionType {
	case "单选题", "多选题", "判断题":
		answerText := strings.TrimSpace(answerSelection.Text())

		// 如果是选项，将中文逗号转换为英文逗号，并清理空格
		answerText = strings.ReplaceAll(answerText, "，", ",")
		answerText = strings.ReplaceAll(answerText, " ", "")
		return strings.Split(answerText, ","), nil

	case "填空题":
		// 填空题：提取每个段落作为答案项
		paragraphs := answerSelection.Find("p")
		var answers []string
		paragraphs.Each(func(i int, s *goquery.Selection) {
			text := cleanHTML(s)
			if text != "" {
				answers = append(answers, text)
			}
		})
		return answers, nil

	case "问答题":
		// 问答题、编程题：保留完整的HTML内容
		return []string{cleanHTML(answerSelection)}, nil

	default:
		// 默认情况：保留完整HTML内容
		return []string{cleanHTML(answerSelection)}, nil
	}
}

// extractAnalysis 提取解析内容
func extractAnalysis(doc *goquery.Document) (string, error) {
	analysisSelection := doc.Find("div.analysis-body")
	if analysisSelection.Length() == 0 {
		return "", fmt.Errorf("未找到解析内容")
	}

	// 如果是文本答案，保留完整的HTML内容（包括图片等）
	return cleanHTML(analysisSelection), nil
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
