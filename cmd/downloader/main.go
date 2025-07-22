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
	"github.com/schollz/progressbar/v3"
)

// 定义图片链接的正则表达式
var imgRegex = regexp.MustCompile(`https?://[^\s<>"]+?/[^\s<>"]+?\.(png|jpg|jpeg|gif|webp)`)

// 定义颜色
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorCyan   = "\033[36m"
)

// 目录信息结构
type DirInfo struct {
	Path      string
	FileCount int
	Files     []string
}

// 获取目录信息
func getDirInfo(root string) (map[string]*DirInfo, error) {
	dirs := make(map[string]*DirInfo)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("遍历目录时出错: %w", err)
		}

		if !d.IsDir() {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				return fmt.Errorf("获取相对路径失败 %s: %w", path, err)
			}

			dir := filepath.Dir(relPath)
			if _, exists := dirs[dir]; !exists {
				dirs[dir] = &DirInfo{
					Path:      dir,
					FileCount: 0,
					Files:     make([]string, 0),
				}
			}
			dirs[dir].FileCount++
			dirs[dir].Files = append(dirs[dir].Files, path)
		}
		return nil
	})

	return dirs, err
}

// 处理单个文件
func processFile(client *api.Client, path string, bar *progressbar.ProgressBar) error {
	// 获取相对于 data 目录的路径
	relPath, err := filepath.Rel("data", path)
	if err != nil {
		return fmt.Errorf("获取相对路径失败 %s: %w", path, err)
	}

	// 构建图片保存目录
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
			bar.Describe(fmt.Sprintf("%s已存在%s %s", colorYellow, colorReset, filename))
		} else {
			// 下载图片
			bar.Describe(fmt.Sprintf("%s正在下载%s %s", colorCyan, colorReset, filename))
			imgData, err := client.Image(imgURL)
			if err != nil {
				return fmt.Errorf("下载图片失败 %s: %w", imgURL, err)
			}

			// 保存图片
			err = os.WriteFile(imgPath, imgData, 0644)
			if err != nil {
				return fmt.Errorf("保存图片失败 %s: %w", imgPath, err)
			}

			bar.Describe(fmt.Sprintf("%s下载完成%s %s", colorGreen, colorReset, filename))
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
		bar.Describe(fmt.Sprintf("%s已更新%s %s", colorBlue, colorReset, filepath.Base(path)))
	}

	return nil
}

func main() {
	// 创建 API 客户端
	client, err := api.NewClient("9c8adef9772544421fc21404af7ed346=e6f5696bd5ff493578540e8afee89af8; ssid=eyJpdiI6IjBRYUlUYzdNeFM4RGhjZ2pOT0RNWnc9PSIsInZhbHVlIjoiZUUzSW5NcTU0bFwvaVYwRFwvaVFIRXhzZldUYWFtcVdVMkpETlZKQmxTb1lkaWNKSGxmT05lVmdNMjA1XC8rQVZhXC84SjY1TG5aWGJRSGlkU0lUblJaZjhBPT0iLCJtYWMiOiI3YWVmYTk0NDQxZjg3YTFiMDlkYjhlNWUwMDNmMzFhM2FkNjBmZWY0Yzg1OWVhYjQ1OTQ0M2E3ODRhNjkzMzI0In0%3D")
	if err != nil {
		fmt.Printf("%s创建客户端失败%s: %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}

	// 获取目录信息
	dirs, err := getDirInfo("data")
	if err != nil {
		fmt.Printf("%s统计目录失败%s: %v\n", colorRed, colorReset, err)
		os.Exit(1)
	}

	// 处理每个目录
	for _, dirInfo := range dirs {
		fmt.Printf("\n%s开始处理目录%s: %s (%d 个文件)\n", colorCyan, colorReset, dirInfo.Path, dirInfo.FileCount)

		// 创建该目录的进度条
		bar := progressbar.NewOptions(dirInfo.FileCount,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(false),
			progressbar.OptionSetWidth(15),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)

		// 处理目录中的每个文件
		for fileIdx, path := range dirInfo.Files {
			// 更新进度条描述
			bar.Describe(fmt.Sprintf("%s正在处理%s (%d/%d) %s", colorCyan, colorReset, fileIdx+1, dirInfo.FileCount, filepath.Base(path)))

			// 处理文件
			if err := processFile(client, path, bar); err != nil {
				fmt.Printf("\n%s错误%s: %v\n", colorRed, colorReset, err)
				os.Exit(1)
			}

			// 更新进度条
			_ = bar.Add(1)
		}

		_ = bar.Finish()
		fmt.Printf("%s目录处理完成%s: %s\n", colorGreen, colorReset, dirInfo.Path)
	}

	fmt.Printf("\n%s所有文件处理完成%s\n", colorGreen, colorReset)
}
