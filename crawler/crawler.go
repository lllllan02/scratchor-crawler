package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"

	"github.com/lllllan02/scratchor-crawler/api"
)

type Crawler struct {
	client *api.Client
	force  bool
}

func NewCrawler(cookie string) (*Crawler, error) {
	client, err := api.NewClient(cookie)
	if err != nil {
		return nil, err
	}

	return &Crawler{client: client}, nil
}

func (c *Crawler) Force() *Crawler {
	c.force = true
	return c
}

func (c *Crawler) Cats() error {
	fmt.Println("开始爬取所有分类...")

	for i := 1; i <= 10; i++ {
		fmt.Printf("正在爬取分类 %d/10...\n", i)
		if err := c.Cat(i); err != nil {
			return err
		}
		fmt.Printf("分类 %d 爬取完成\n", i)
	}

	fmt.Println("所有分类爬取完成")
	return nil
}

func (c *Crawler) Cat(id int) error {
	fmt.Printf("开始爬取分类 %d 的章节...\n", id)

	chatpters, err := c.client.GetCat(id)
	if err != nil {
		fmt.Printf("failed crawl cat %d: %v\n", id, err)
		return err
	}

	fmt.Printf("分类 %d 共有 %d 个章节\n", id, len(chatpters))

	for i, chapter := range chatpters {
		fmt.Printf("正在爬取分类 %d 的章节 %d/%d: %s\n", id, i+1, len(chatpters), chapter)
		if err := c.Chapter(chapter); err != nil {
			return err
		}
		fmt.Printf("分类 %d 的章节 %d/%d 爬取完成\n", id, i+1, len(chatpters))
	}

	fmt.Printf("分类 %d 所有章节爬取完成\n", id)

	return nil
}

func (c *Crawler) Chapter(url string) error {
	// 从 URL 中提取 cat, chapter
	var cat, chapter int
	if _, err := fmt.Sscanf(url, "https://tiku.scratchor.com/question/cat/%d/list?chapterId=%d", &cat, &chapter); err != nil {
		return fmt.Errorf("解析 chapter ID 失败: %v", err)
	}
	path := fmt.Sprintf("data/cat_%d/chapter_%d", cat, chapter)

	fmt.Printf("开始爬取分类 %d 章节 %d，保存路径: %s\n", cat, chapter, path)

	// 创建目录
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %v", err)
	}
	fmt.Printf("目录创建成功: %s\n", path)

	pageNum := 1
	for {
		fmt.Printf("正在爬取分类 %d 章节 %d 第 %d 页...\n", cat, chapter, pageNum)

		views, next, err := c.client.GetChapter(url)
		if err != nil {
			fmt.Printf("failed crawl chapter %s: %v\n", url, err)
			return err
		}

		fmt.Printf("分类 %d 章节 %d 第 %d 页获取到 %d 个题目\n", cat, chapter, pageNum, len(views))

		for i, view := range views {
			fmt.Printf("正在处理分类 %d 章节 %d 第 %d 页的题目 %d/%d: %s\n", cat, chapter, pageNum, i+1, len(views), view)
			if err := c.View(path, view); err != nil {
				return err
			}
		}

		if next == "" {
			fmt.Printf("分类 %d 章节 %d 所有页面爬取完成，共 %d 页\n", cat, chapter, pageNum)
			break
		}
		url = next
		pageNum++
	}

	fmt.Printf("分类 %d 章节 %d 爬取完成\n", cat, chapter)

	return nil
}

func (c *Crawler) View(path, url string) error {
	// 从 URL 中提取 alias
	var alias string
	if _, err := fmt.Sscanf(url, "https://tiku.scratchor.com/question/view/%s", &alias); err != nil {
		return fmt.Errorf("解析 URL 中的 alias 失败: %v", err)
	}

	// 构建文件路径
	filePath := fmt.Sprintf("%s/%s.json", path, alias)

	// 如果不是强制模式且文件已存在，则跳过
	if !c.force {
		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("文件已存在，跳过: %s\n", filePath)
			return nil
		}
	}

	fmt.Printf("开始获取题目详情: %s (alias: %s)\n", url, alias)

	view, err := c.client.GetView(url)
	if err != nil {
		fmt.Printf("failed crawl view %s: %v\n", url, err)
		return err
	}

	fmt.Printf("题目详情获取成功，别名: %s，保存路径: %s\n", alias, filePath)

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

	fmt.Printf("题目详情保存成功: %s (大小: %d 字节)\n", filePath, len(data))

	return nil
}
