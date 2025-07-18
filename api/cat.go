package api

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func (client *Client) GetCat(id int) ([]string, error) {
	url := fmt.Sprintf("https://tiku.scratchor.com/question/cat/%d", id)

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

	// 找出所有章节链接
	var links []string
	doc.Find("a.ub-text-primary").Each(func(i int, s *goquery.Selection) {
		if href, exists := s.Attr("href"); exists {
			links = append(links, fmt.Sprintf("https://tiku.scratchor.com%s", href))
		}
	})

	return links, nil
}
