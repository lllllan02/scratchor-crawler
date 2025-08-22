package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lllllan02/scratchor-crawler/api"
	"github.com/lllllan02/scratchor-crawler/utils"
	"github.com/schollz/progressbar/v3"
)

var client *api.Client

func main() {
	var err error
	if client, err = api.NewClient("9c8adef9772544421fc21404af7ed346=e6f5696bd5ff493578540e8afee89af8; ssid=eyJpdiI6IjBRYUlUYzdNeFM4RGhjZ2pOT0RNWnc9PSIsInZhbHVlIjoiZUUzSW5NcTU0bFwvaVYwRFwvaVFIRXhzZldUYWFtcVdVMkpETlZKQmxTb1lkaWNKSGxmT05lVmdNMjA1XC8rQVZhXC84SjY1TG5aWGJRSGlkU0lUblJaZjhBPT0iLCJtYWMiOiI3YWVmYTk0NDQxZjg3YTFiMDlkYjhlNWUwMDNmMzFhM2FkNjBmZWY0Yzg1OWVhYjQ1OTQ0M2E3ODRhNjkzMzI0In0%3D"); err != nil {
		panic(err)
	}

	Cats()
}

func Cats() error {
	startID := 1
	if maxID, err := findMaxCatID(); err == nil {
		startID = maxID
	}

	// 创建分类爬取进度条
	totalCats := 10 - startID + 1
	catBar := progressbar.NewOptions(totalCats,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("%s爬取分类进度%s", utils.ColorCyan, utils.ColorReset)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for i := startID; i <= 10; i++ {
		catBar.Describe(fmt.Sprintf("%s正在爬取分类 %d/10%s", utils.ColorBlue, i, utils.ColorReset))

		if err := Cat(i); err != nil {
			catBar.Describe(fmt.Sprintf("%s分类 %d 爬取失败%s", utils.ColorRed, i, utils.ColorReset))
			return err
		}

		catBar.Describe(fmt.Sprintf("%s分类 %d 爬取完成%s", utils.ColorGreen, i, utils.ColorReset))
		_ = catBar.Add(1)
	}

	_ = catBar.Finish()
	return nil
}

// 查找已存在的最大分类ID
func findMaxCatID() (int, error) {
	// 检查data目录是否存在
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		return 0, fmt.Errorf("data目录不存在")
	}

	maxID := 0
	entries, err := os.ReadDir("data")
	if err != nil {
		return 0, fmt.Errorf("读取data目录失败: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 解析目录名中的ID
		var id int
		if _, err := fmt.Sscanf(entry.Name(), "cat_%d", &id); err != nil {
			continue
		}

		if id > maxID {
			maxID = id
		}
	}

	if maxID == 0 {
		return 0, fmt.Errorf("未找到有效的分类目录")
	}

	return maxID, nil
}

func Cat(id int) error {
	links, err := client.GetCat(id)
	if err != nil {
		fmt.Printf("%sfailed crawl cat %d: %v%s\n", utils.ColorRed, id, err, utils.ColorReset)
		return err
	}

	// 提取链接列表
	var chatpters []string
	for _, link := range links {
		chatpters = append(chatpters, link.URL)
	}

	startIdx := 0
	if maxChapter, err := findMaxChapter(id); err == nil {
		// 找到对应的索引
		for i, chapter := range chatpters {
			var cat, chap int
			if _, err := fmt.Sscanf(chapter, "https://tiku.scratchor.com/question/cat/%d/list?chapterId=%d", &cat, &chap); err == nil {
				if chap == maxChapter {
					startIdx = i
					break
				}
			}
		}
	}

	// 创建章节爬取进度条
	totalChapters := len(chatpters) - startIdx
	chapterBar := progressbar.NewOptions(totalChapters,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15),
		progressbar.OptionSetDescription(fmt.Sprintf("%s分类 %d 章节爬取进度%s", utils.ColorCyan, id, utils.ColorReset)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for i := startIdx; i < len(chatpters); i++ {
		chapterBar.Describe(fmt.Sprintf("%s正在爬取分类 %d 的章节 %d/%d: %s%s", utils.ColorBlue, id, i+1, len(chatpters), chatpters[i], utils.ColorReset))

		if err := Chapter(chatpters[i]); err != nil {
			chapterBar.Describe(fmt.Sprintf("%s分类 %d 的章节 %d/%d 爬取失败%s", utils.ColorRed, id, i+1, len(chatpters), utils.ColorReset))
			return err
		}

		chapterBar.Describe(fmt.Sprintf("%s分类 %d 的章节 %d/%d 爬取完成%s", utils.ColorGreen, id, i+1, len(chatpters), utils.ColorReset))
		_ = chapterBar.Add(1)
	}

	_ = chapterBar.Finish()

	return nil
}

// 查找已存在的最大章节ID
func findMaxChapter(catID int) (int, error) {
	// 检查分类目录是否存在
	path := fmt.Sprintf("data/cat_%d", catID)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return 0, fmt.Errorf("分类目录不存在")
	}

	// 读取目录下的所有文件夹
	entries, err := os.ReadDir(path)
	if err != nil {
		return 0, fmt.Errorf("读取目录失败: %w", err)
	}

	maxChapter := 0
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		// 解析目录名中的chapter ID
		var chapter int
		if _, err := fmt.Sscanf(entry.Name(), "chapter_%d", &chapter); err != nil {
			continue
		}

		if chapter > maxChapter {
			maxChapter = chapter
		}
	}

	if maxChapter == 0 {
		return 0, fmt.Errorf("未找到有效的章节目录")
	}

	return maxChapter, nil
}

func Chapter(url string) error {
	// 从 URL 中提取 cat, chapter
	var cat, chapter int
	if _, err := fmt.Sscanf(url, "https://tiku.scratchor.com/question/cat/%d/list?chapterId=%d", &cat, &chapter); err != nil {
		return fmt.Errorf("解析 chapter ID 失败: %v", err)
	}
	path := fmt.Sprintf("data/cat_%d/chapter_%d", cat, chapter)

	// 创建目录
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}

	pageNum := 1
	totalViews := 0

	for {
		views, next, err := client.GetChapter(url)
		if err != nil {
			fmt.Printf("%sfailed crawl chapter %s: %v%s\n", utils.ColorRed, url, err, utils.ColorReset)
			return err
		}

		// 创建题目处理进度条
		if len(views) > 0 {
			viewBar := progressbar.NewOptions(len(views),
				progressbar.OptionEnableColorCodes(true),
				progressbar.OptionShowBytes(false),
				progressbar.OptionSetWidth(15),
				progressbar.OptionSetDescription(fmt.Sprintf("%s第 %d 页题目处理进度%s", utils.ColorCyan, pageNum, utils.ColorReset)),
				progressbar.OptionSetTheme(progressbar.Theme{
					Saucer:        "[green]=[reset]",
					SaucerHead:    "[green]>[reset]",
					SaucerPadding: " ",
					BarStart:      "[",
					BarEnd:        "]",
				}),
			)

			for i, view := range views {
				viewBar.Describe(fmt.Sprintf("%s正在处理题目 %d/%d: %s%s", utils.ColorBlue, i+1, len(views), view, utils.ColorReset))

				if err := View(path, view); err != nil {
					viewBar.Describe(fmt.Sprintf("%s题目 %d/%d 处理失败%s", utils.ColorRed, i+1, len(views), utils.ColorReset))
					return err
				}

				viewBar.Describe(fmt.Sprintf("%s题目 %d/%d 处理完成%s", utils.ColorGreen, i+1, len(views), utils.ColorReset))
				_ = viewBar.Add(1)
			}

			_ = viewBar.Finish()
		}

		totalViews += len(views)

		if next == "" {
			break
		}
		url = next
		pageNum++
	}

	return nil
}

func View(path, url string) error {
	// 从 URL 中提取 alias
	var alias string
	if _, err := fmt.Sscanf(url, "https://tiku.scratchor.com/question/view/%s", &alias); err != nil {
		return fmt.Errorf("解析 URL 中的 alias 失败: %v", err)
	}

	// 构建文件路径
	filePath := fmt.Sprintf("%s/%s.json", path, alias)

	if _, err := os.Stat(filePath); err == nil {
		return nil
	}

	view, err := client.GetView(url)
	if err != nil {
		fmt.Printf("%sfailed crawl view %s: %v%s\n", utils.ColorRed, url, err, utils.ColorReset)
		return err
	}

	// 将 view 转换为 JSON，不对字符进行转义
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(view); err != nil {
		return fmt.Errorf("JSON 序列化失败: %v", err)
	}
	data := buf.Bytes()

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}
