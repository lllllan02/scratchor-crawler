package api

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GetChapter(url, cookie string) (links []string, next string, err error) {
	// 获取 html
	body, err := Get(url, cookie)
	if err != nil {
		return nil, "", err
	}

	// 解析 html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, "", err
	}

	// 找出所有题目链接
	doc.Find("a.question-title").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			links = append(links, fmt.Sprintf("https://tiku.scratchor.com%s", href))
		}
	})

	// 找到「下一页」的链接
	doc.Find("a.page:contains('下一页')").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			baseURL := strings.Split(url, "?")[0]
			next = fmt.Sprintf("%s%s", baseURL, href)
		}
	})

	return links, next, nil
}
