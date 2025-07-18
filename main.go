package main

import "github.com/lllllan02/scratchor-crawler/crawler"

func main() {
	craler, err := crawler.NewCrawler("26efeb6279834e831e8659742d83d367=04310906d8405e72b302c5514f0d1e0b; ssid=eyJpdiI6IkNpU0lhYWpmclEyNTVOb3dmZ1JrZUE9PSIsInZhbHVlIjoiY3JFeW5oVkVKMWw3U0tQcnRzd1hJUTQ0Y1JRMHZ1R2ZKYno4bGl3XC83MWRjRmNKZzF0U1JQSU5wbWxLWVh0eEJlcGVqOXp1XC9KaVYwOGNNdElZZHBoZz09IiwibWFjIjoiZGM0MGIyN2RlYmIwMjQ5ODEwNjNkZTJlNWU1ZDgzYzdmYThjNGVmM2E4OTQxOTQ2YTQ1ZTQxYjNiMWU3Y2MzOSJ9")

	if err != nil {
		panic(err)
	}

	craler.Cats()
}
