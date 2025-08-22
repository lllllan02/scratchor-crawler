package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/lllllan02/scratchor-crawler/api"
)

// 图片下载器配置
type ImageDownloaderConfig struct {
	DataDir string      // 数据目录
	ImgDir  string      // 图片保存目录
	Client  *api.Client // API客户端
}

// 图片链接正则表达式
var ImgRegex = regexp.MustCompile(`https?://[^\s<>"]+?/[^\s<>"]+?\.(png|jpg|jpeg|gif|webp)`)

// 创建图片下载器处理函数
func CreateImageDownloader(config ImageDownloaderConfig) func(string) (bool, error) {
	return func(filePath string) (bool, error) {
		// 获取相对于 data 目录的路径
		relPath, err := filepath.Rel(config.DataDir, filePath)
		if err != nil {
			return false, fmt.Errorf("获取相对路径失败 %s: %w", filePath, err)
		}

		// 构建图片保存目录
		imgDir := filepath.Join(config.ImgDir, strings.TrimSuffix(relPath, filepath.Ext(relPath)))
		if err := EnsureDir(imgDir); err != nil {
			return false, fmt.Errorf("创建目录失败 %s: %w", imgDir, err)
		}

		// 读取文件内容
		content, err := ReadFile(filePath)
		if err != nil {
			return false, fmt.Errorf("读取文件失败 %s: %w", filePath, err)
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
		matches := ImgRegex.FindAllString(contentStr, -1)
		for _, match := range matches {
			// 清理链接
			imgURL := strings.TrimSpace(match)

			// 从 URL 中提取文件名
			filename := filepath.Base(imgURL)
			imgPath := filepath.Join(imgDir, filename)

			// 检查文件是否已存在
			if FileExists(imgPath) {
				// 图片已存在，跳过下载
				continue
			}

			// 下载图片
			imgData, err := config.Client.Image(imgURL)
			if err != nil {
				return false, fmt.Errorf("下载图片失败 %s: %w", imgURL, err)
			}

			// 保存图片
			err = WriteFile(imgPath, imgData)
			if err != nil {
				return false, fmt.Errorf("保存图片失败 %s: %w", imgPath, err)
			}

			// 计算图片的相对路径（相对于 JSON 文件）
			relImgPath := filepath.Join("..", config.ImgDir, strings.TrimSuffix(relPath, filepath.Ext(relPath)), filename)

			// 替换文件中的图片链接
			if strings.Contains(contentStr, imgURL) {
				contentStr = strings.ReplaceAll(contentStr, imgURL, relImgPath)
				needUpdate = true
			}
		}

		// 如果有更新，保存文件
		if needUpdate {
			var newContent []byte
			var err error
			if isJSON {
				var buf bytes.Buffer
				encoder := json.NewEncoder(&buf)
				encoder.SetEscapeHTML(false)
				encoder.SetIndent("", "  ")
				if err := encoder.Encode(jsonData); err != nil {
					return false, fmt.Errorf("序列化 JSON 失败 %s: %w", filePath, err)
				}
				newContent = buf.Bytes()
			} else {
				newContent = []byte(contentStr)
			}

			err = WriteFile(filePath, newContent)
			if err != nil {
				return false, fmt.Errorf("更新文件失败 %s: %w", filePath, err)
			}
		}

		return needUpdate, nil
	}
}
