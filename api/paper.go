package api

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (client *Client) GetPaper(url string) (links []string, next string, err error) {
	// 获取 html
	body, err := client.Get(url)
	if err != nil {
		return nil, "", err
	}

	// 解析 html
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(body))
	if err != nil {
		return nil, "", err
	}

	// 找出所有试卷链接
	doc.Find("a.btn-round").Each(func(i int, s *goquery.Selection) {
		if paper, exist := s.Attr("data-paper-view"); exist {
			links = append(links, fmt.Sprintf("https://tiku.scratchor.com/paper/view/%s", paper))
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
