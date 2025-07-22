package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lllllan02/scratchor-crawler/api"
)

// 定义图片链接的正则表达式
var imgRegex = regexp.MustCompile(`https?://[^\s<>"]+?/[^\s<>"]+?\.(png|jpg|jpeg|gif|webp)`)

func main() {
	// 创建 API 客户端
	client, err := api.NewClient("9c8adef9772544421fc21404af7ed346=e6f5696bd5ff493578540e8afee89af8; ssid=eyJpdiI6IjBRYUlUYzdNeFM4RGhjZ2pOT0RNWnc9PSIsInZhbHVlIjoiZUUzSW5NcTU0bFwvaVYwRFwvaVFIRXhzZldUYWFtcVdVMkpETlZKQmxTb1lkaWNKSGxmT05lVmdNMjA1XC8rQVZhXC84SjY1TG5aWGJRSGlkU0lUblJaZjhBPT0iLCJtYWMiOiI3YWVmYTk0NDQxZjg3YTFiMDlkYjhlNWUwMDNmMzFhM2FkNjBmZWY0Yzg1OWVhYjQ1OTQ0M2E3ODRhNjkzMzI0In0%3D")
	if err != nil {
		fmt.Printf("创建客户端失败: %v\n", err)
		os.Exit(1)
	}

	// 遍历 data 目录
	err = filepath.WalkDir("data", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("遍历目录时出错: %w", err)
		}

		// 跳过目录
		if d.IsDir() {
			return nil
		}

		// 获取相对于 data 目录的路径
		relPath, err := filepath.Rel("data", path)
		if err != nil {
			return fmt.Errorf("获取相对路径失败 %s: %w", path, err)
		}

		// 构建图片保存目录
		// 例如：data/cat_1/chapter_1/abc.json -> image/cat_1/chapter_1/abc/
		imgDir := filepath.Join("image", strings.TrimSuffix(relPath, filepath.Ext(relPath)))
		if err := os.MkdirAll(imgDir, 0755); err != nil {
			return fmt.Errorf("创建目录失败 %s: %w", imgDir, err)
		}

		// 读取文件内容
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("读取文件失败 %s: %w", path, err)
		}

		// 记录是否为 JSON 文件
		isJSON := false
		var jsonData interface{}

		// 尝试解析 JSON
		if err := json.Unmarshal(content, &jsonData); err == nil {
			isJSON = true
		}

		// 记录是否需要更新文件
		needUpdate := false
		contentStr := string(content)

		// 查找所有图片链接
		matches := imgRegex.FindAllString(contentStr, -1)
		for _, match := range matches {
			// 清理链接
			imgURL := strings.TrimSpace(match)

			// 从 URL 中提取文件名
			filename := filepath.Base(imgURL)
			imgPath := filepath.Join(imgDir, filename)

			// 检查文件是否已存在
			if _, err := os.Stat(imgPath); err == nil {
				fmt.Printf("图片已存在，跳过: %s\n", imgPath)
			} else {
				// 下载图片
				fmt.Printf("正在下载: %s\n", imgURL)
				imgData, err := client.Image(imgURL)
				if err != nil {
					return fmt.Errorf("下载图片失败 %s: %w", imgURL, err)
				}

				// 保存图片
				err = os.WriteFile(imgPath, imgData, 0644)
				if err != nil {
					return fmt.Errorf("保存图片失败 %s: %w", imgPath, err)
				}

				fmt.Printf("已保存: %s\n", imgPath)
			}

			// 计算图片的相对路径（相对于 JSON 文件）
			relImgPath := filepath.Join("..", "image", strings.TrimSuffix(relPath, filepath.Ext(relPath)), filename)
			if isJSON {
				// 如果是 JSON 文件，直接在原始内容中替换
				if strings.Contains(contentStr, imgURL) {
					contentStr = strings.ReplaceAll(contentStr, imgURL, relImgPath)
					needUpdate = true
				}
			} else {
				// 如果是普通文件，直接替换内容
				if strings.Contains(contentStr, imgURL) {
					contentStr = strings.ReplaceAll(contentStr, imgURL, relImgPath)
					needUpdate = true
				}
			}
		}

		// 如果有更新，保存文件
		if needUpdate {
			var newContent []byte
			if isJSON {
				// 如果是 JSON 文件，先解析再格式化
				var jsonObj interface{}
				if err := json.Unmarshal([]byte(contentStr), &jsonObj); err != nil {
					return fmt.Errorf("解析更新后的 JSON 失败 %s: %w", path, err)
				}
				newContent, err = json.MarshalIndent(jsonObj, "", "  ")
				if err != nil {
					return fmt.Errorf("序列化 JSON 失败 %s: %w", path, err)
				}
			} else {
				newContent = []byte(contentStr)
			}

			err = os.WriteFile(path, newContent, 0644)
			if err != nil {
				return fmt.Errorf("更新文件失败 %s: %w", path, err)
			}
			fmt.Printf("已更新文件: %s\n", path)
		}

		return nil
	})

	if err != nil {
		fmt.Printf("错误: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("所有文件处理完成")
}
