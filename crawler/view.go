package crawler

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type View struct {
	*Question
	Tags  []string
	Items []*Question
}

type Question struct {
	Alias  string
	Type   string
	Body   string
	Option []string
}

func GetView(url, cookie string) (*View, error) {
	body, err := Get(url, cookie)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	view := &View{}
	// 获取题目
	view.Question = getQuestion(doc.Find("div.pb-question-view"))

	// 获取题目标签
	doc.Find("div.ub-panel .body span").Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			view.Tags = append(view.Tags, text)
		}
	})

	// 获取子题目
	doc.Find("div.question-items .pb-question-view").Each(func(i int, s *goquery.Selection) {
		view.Items = append(view.Items, getQuestion(s))
	})

	return view, nil
}

func getQuestion(doc *goquery.Selection) *Question {
	question := &Question{}

	// 从 div.question 中获取 alias
	if alias, exists := doc.Find("div.question").Attr("data-question-alias"); exists {
		question.Alias = alias
	}

	// 获取题目类型
	doc.Find("div.pb-question-view .title").Each(func(i int, s *goquery.Selection) {
		question.Type = getType(s.Text())
	})

	// 获取题目内容
	doc.Find("div.question .ub-html").Each(func(i int, s *goquery.Selection) {
		question.Body, _ = s.Html()
		question.Body = strings.TrimSpace(question.Body)
	})

	// 获取选项
	doc.Find("div.question .option .item").Each(func(i int, s *goquery.Selection) {
		html, _ := s.Find("p").Html()
		question.Option = append(question.Option, strings.TrimSpace(html))
	})

	return question
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
