package api

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Question struct {
	*QuestionBody
	Tags  []string
	Items []*QuestionBody
}

type QuestionBody struct {
	Alias    string   // 题目别名
	Type     string   // 题目类型
	Body     string   // 题目内容
	Analysis string   // 题目解析
	Option   []string // 题目选项
	Answer   []string // 题目答案
}

func (client *Client) GetQuestion(url string) (*Question, error) {
	body, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	view := &Question{}

	// 先获取子题目（在移除之前）
	doc.Find("div.question-items .pb-question-view").Each(func(i int, s *goquery.Selection) {
		view.Items = append(view.Items, getQuestion(s))
	})

	// 获取主干题目（排除子题部分）
	mainQuestionDoc := doc.Find("div.pb-question-view").First()
	// 移除子题部分，确保主干题目不包含子题内容
	mainQuestionDoc.Find("div.question-items").Remove()
	view.QuestionBody = getQuestion(mainQuestionDoc)

	// 获取题目标签
	doc.Find("div.ub-panel .body span").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			view.Tags = append(view.Tags, text)
		}
	})

	return view, nil
}

func getQuestion(doc *goquery.Selection) *QuestionBody {
	question := &QuestionBody{}

	// 从 div.question 中获取 alias
	if alias, exists := doc.Find("div.question").Attr("data-question-alias"); exists {
		question.Alias = alias
	}
	if alias, exists := doc.Find("div.body").Attr("data-question-alias"); exists {
		question.Alias = alias
	}

	// 获取题目类型
	doc.Find("div.pb-question-view .title").Each(func(i int, s *goquery.Selection) {
		question.Type = getType(s.Text())
	})

	// 获取题目内容
	question.Body = cleanHTML(doc.Find("div.question .ub-html").First())

	// 获取选项
	doc.Find("div.question .option .item").Each(func(i int, s *goquery.Selection) {
		// 移除选择器元素
		cleanOption := s.Clone()
		cleanOption.Find(".choice").Remove()
		question.Option = append(question.Option, cleanHTML(cleanOption))
	})

	return question
}

// cleanHTML 清理HTML内容，去除样式标签和属性
func cleanHTML(selection *goquery.Selection) string {
	// 克隆选择器，避免修改原始数据
	clean := selection.Clone()

	// 移除样式标签
	clean.Find("style, script, link, meta").Remove()

	// 移除所有元素的样式相关属性
	clean.Find("*").Each(func(i int, s *goquery.Selection) {
		s.RemoveAttr("width")
		s.RemoveAttr("height")
		s.RemoveAttr("style")
		s.RemoveAttr("class")
		s.RemoveAttr("id")
	})

	// 获取清理后的HTML
	html, _ := clean.Html()
	return strings.TrimSpace(html)
}

func getType(t string) string {
	if strings.Contains(t, "单选") {
		return "单选"
	}
	if strings.Contains(t, "多选") {
		return "多选"
	}
	if strings.Contains(t, "判断") {
		return "判断"
	}
	if strings.Contains(t, "填空") {
		return "填空"
	}
	if strings.Contains(t, "问答") {
		return "问答"
	}
	if strings.Contains(t, "组合") {
		return "组合"
	}
	return "未知"
}
