package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lllllan02/scratchor-crawler/api"
	"github.com/lllllan02/scratchor-crawler/utils"
	"github.com/schollz/progressbar/v3"
)

var client *api.Client

func main() {
	// 解析命令行参数
	var limiterConcurrency uint
	flag.UintVar(&limiterConcurrency, "limiter", 2, "限流器并发数，默认为2")
	flag.Parse()

	// 设置全局限流器并发数
	api.LimiterConcurrency = limiterConcurrency

	var err error
	if client, err = api.NewClient("9c8adef9772544421fc21404af7ed346=e6f5696bd5ff493578540e8afee89af8; ssid=eyJpdiI6IjBRYUlUYzdNeFM4RGhjZ2pOT0RNWnc9PSIsInZhbHVlIjoiZUUzSW5NcTU0bFwvaVYwRFwvaVFIRXhzZldUYWFtcVdVMkpETlZKQmxTb1lkaWNKSGxmT05lVmdNMjA1XC8rQVZhXC84SjY1TG5aWGJRSGlkU0lUblJaZjhBPT0iLCJtYWMiOiI3YWVmYTk0NDQxZjg3YTFiMDlkYjhlNWUwMDNmMzFhM2FkNjBmZWY0Yzg1OWVhYjQ1OTQ0M2E3ODRhNjkzMzI0In0%3D"); err != nil {
		panic(err)
	}

	Cats()
}

func Cats() error {
	for i := 1; i <= 10; i++ {
		if err := Cat(i); err != nil {
			return err
		}
	}
	return nil
}

func Cat(id int) error {
	cat, err := client.GetCat(id)
	if err != nil {
		fmt.Printf("%sfailed crawl cat %d: %v%s\n", utils.ColorRed, id, err, utils.ColorReset)
		return err
	}

	fmt.Printf("\n%s获取专题「%s」: %d 个章节%s\n", utils.ColorCyan, cat.Title, len(cat.Links), utils.ColorReset)
	for i, link := range cat.Links {
		bar := progressbar.NewOptions(
			link.Count,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(false),
			progressbar.OptionSetWidth(50),
			progressbar.OptionSetDescription(fmt.Sprintf("%s获取练习「%s」%d/%d: %d 个题目%s", utils.ColorCyan, link.Title, i+1, len(cat.Links), link.Count, utils.ColorReset)),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)

		path := fmt.Sprintf("data/%s/%s", cat.Title, link.Title)
		if err := Chapter(path, link.URL, bar); err != nil {
			return err
		}

		bar.Finish()
		fmt.Println()
	}

	return nil
}

func Chapter(path, url string, bar *progressbar.ProgressBar) error {
	// 从 URL 中提取 cat, chapter
	var cat, chapter int
	if _, err := fmt.Sscanf(url, "https://tiku.scratchor.com/question/cat/%d/list?chapterId=%d", &cat, &chapter); err != nil {
		return fmt.Errorf("解析 chapter ID 失败: %v", err)
	}

	// 创建目录
	os.MkdirAll(path, 0755)

	for pageNum := 1; url != ""; {
		views, next, err := client.GetChapter(url)
		if err != nil {
			fmt.Printf("%s获取练习「%s」失败: %v%s\n", utils.ColorRed, url, err, utils.ColorReset)
			return err
		}

		// 创建题目处理进度条
		for _, view := range views {
			progressbar.Bprintln(bar, fmt.Sprintf("获取题目「%s」", view))

			if err := View(path, view); err != nil {
				return err
			}

			bar.Add(1)
		}

		url, pageNum = next, pageNum+1
	}

	return nil
}

func View(path, url string) error {
	// 从 URL 中提取 alias
	var alias string
	if _, err := fmt.Sscanf(url, "https://tiku.scratchor.com/question/view/%s", &alias); err != nil {
		fmt.Printf("%s解析 URL 中的 alias 失败: %v%s\n", utils.ColorRed, err, utils.ColorReset)
		return err
	}

	// 构建文件路径
	filePath := fmt.Sprintf("%s/%s.json", path, alias)

	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	view, err := client.GetView(url)
	if err != nil {
		fmt.Printf("%s获取题目「%s」失败: %v%s\n", utils.ColorRed, url, err, utils.ColorReset)
		return err
	}

	// 将 view 转换为 JSON 并写入文件
	if err := utils.WriteJSON(filePath, view); err != nil {
		fmt.Printf("%s%s%s\n", utils.ColorRed, err, utils.ColorReset)
		return err
	}

	return nil
}
