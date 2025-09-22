package api

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type Paper struct {
	Head           string
	TotalQuestions string
	TotalScore     string
	Duration       string
	Sections       []*Section
}

type Section struct {
	Head      string
	Questions []string
}

func (client *Client) GetPaper(url string) (*Paper, error) {
	// 获取 html
	body, err := client.Get(url)
	if err != nil {
		return nil, err
	}

	// 解析 html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, err
	}

	paper := &Paper{}

	// 获取试卷标题
	title := doc.Find("div.ub-panel div.title").First()
	paper.Head = strings.TrimSpace(title.Text())

	// 获取试卷信息
	doc.Find("div.pb-paper-exam-summary div.line").Each(func(i int, s *goquery.Selection) {
		if strings.Contains(s.Text(), "题目总数") {
			paper.TotalQuestions = strings.TrimSpace(strings.Split(s.Text(), "：")[1])
		} else if strings.Contains(s.Text(), "总分数") {
			paper.TotalScore = strings.TrimSpace(strings.Split(s.Text(), "：")[1])
		} else if strings.Contains(s.Text(), "时间") {
			paper.Duration = strings.TrimSpace(strings.Split(s.Text(), "：")[1])
		}
	})

	doc.Find("div.section").Each(func(i int, s *goquery.Selection) {
		section := &Section{}

		head := s.Find("div.section-head").First()
		section.Head = parseSectionHead(head.Text())

		s.Find("div.section-body a").Each(func(i int, a *goquery.Selection) {
			if alias, exit := a.Attr("data-alias"); exit {
				section.Questions = append(section.Questions, strings.TrimRight(alias, "-0"))
			}
		})

		paper.Sections = append(paper.Sections, section)
	})

	return paper, nil
}

func parseSectionHead(text string) string {
	if idx := strings.Index(text, "、"); idx >= 0 {
		return strings.TrimSpace(text[idx+len("、"):])
	}
	return strings.TrimSpace(text)
}
