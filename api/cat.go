package api

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type LinkInfo struct {
	URL   string
	Count int
}

func (client *Client) GetCat(id int) ([]*LinkInfo, error) {
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

	var links []*LinkInfo

	panel := doc.Find("div.ub-panel").First()
	panel.Find("div.tw-flex").Each(func(i int, row *goquery.Selection) {
		link := &LinkInfo{}

		text := row.Find("div.ub-text-muted").First().Text()
		fmt.Sscanf(strings.TrimSpace(text), "%d题", &link.Count)
		if href, exist := row.Find("a.ub-text-primary").First().Attr("href"); exist {
			link.URL = fmt.Sprintf("https://tiku.scratchor.com%s", href)
		}

		links = append(links, link)
	})

	return links, nil
}
