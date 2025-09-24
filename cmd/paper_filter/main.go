package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// PaperMapItem 表示 paper_map.json 中的项目结构
type PaperMapItem struct {
	Alias    string         `json:"Alias"`
	Url      string         `json:"Url"`
	Title    string         `json:"Title"`
	Children []PaperMapItem `json:"Children"`
}

// ClassificationResult 表示分类结果
type ClassificationResult struct {
	FilePath string
	Category string
	Level    string
	Type     string
}

func main() {
	// 读取 paper_map.json
	paperMapData, err := ioutil.ReadFile("paper_map.json")
	if err != nil {
		fmt.Printf("读取 paper_map.json 失败: %v\n", err)
		return
	}

	var paperMap []PaperMapItem
	err = json.Unmarshal(paperMapData, &paperMap)
	if err != nil {
		fmt.Printf("解析 paper_map.json 失败: %v\n", err)
		return
	}

	// 创建 alias 到分类路径的映射
	aliasMap := make(map[string]string)
	buildAliasMap(paperMap, "", aliasMap)

	// 读取 paper 目录下的所有文件
	paperDir := "paper"
	files, err := ioutil.ReadDir(paperDir)
	if err != nil {
		fmt.Printf("读取 paper 目录失败: %v\n", err)
		return
	}

	// 分类结果
	var classifications []ClassificationResult
	var unclassified []string

	// 对每个文件进行分类
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		// 去掉 .json 扩展名获取 alias
		alias := strings.TrimSuffix(file.Name(), ".json")

		if categoryPath, exists := aliasMap[alias]; exists {
			classifications = append(classifications, ClassificationResult{
				FilePath: file.Name(),
				Category: categoryPath,
			})
		} else {
			unclassified = append(unclassified, file.Name())
		}
	}

	// 输出分类结果
	fmt.Printf("=== 文件分类结果 ===\n")
	fmt.Printf("总文件数: %d\n", len(files))
	fmt.Printf("已分类: %d\n", len(classifications))
	fmt.Printf("未分类: %d\n", len(unclassified))

	// 按分类路径分组显示
	categoryGroups := make(map[string][]string)
	for _, result := range classifications {
		categoryGroups[result.Category] = append(categoryGroups[result.Category], result.FilePath)
	}

	fmt.Printf("\n=== 分类详情 ===\n")
	for category, files := range categoryGroups {
		fmt.Printf("\n%s (%d 个文件):\n", category, len(files))
		for _, file := range files {
			fmt.Printf("  - %s\n", file)
		}
	}

	if len(unclassified) > 0 {
		fmt.Printf("\n=== 未分类文件 ===\n")
		for _, file := range unclassified {
			fmt.Printf("  - %s\n", file)
		}
	}

	// 创建分类目录结构
	createDirectoryStructure(categoryGroups)

	// 复制文件到对应目录
	copyFilesToDirectories(classifications)
}

// buildAliasMap 递归构建 alias 到分类路径的映射
func buildAliasMap(items []PaperMapItem, parentPath string, aliasMap map[string]string) {
	for _, item := range items {
		currentPath := parentPath
		if item.Title != "" {
			// 只替换 Title 中的 / 为 _，保持目录层次结构
			safeTitle := strings.ReplaceAll(item.Title, "/", "_")
			if currentPath != "" {
				currentPath += "/" + safeTitle
			} else {
				currentPath = safeTitle
			}
		}

		if item.Alias != "" {
			aliasMap[item.Alias] = currentPath
		}

		if item.Children != nil {
			buildAliasMap(item.Children, currentPath, aliasMap)
		}
	}
}

// createDirectoryStructure 创建分类目录结构
func createDirectoryStructure(categoryGroups map[string][]string) {
	baseDir := "paper2"

	for category := range categoryGroups {
		dirPath := filepath.Join(baseDir, category)
		err := os.MkdirAll(dirPath, 0755)
		if err != nil {
			fmt.Printf("创建目录失败 %s: %v\n", dirPath, err)
			continue
		}
		fmt.Printf("创建目录: %s\n", dirPath)
	}
}

// copyFilesToDirectories 将文件复制到对应的分类目录
func copyFilesToDirectories(classifications []ClassificationResult) {
	baseDir := "paper2"
	sourceDir := "paper"

	copiedCount := 0
	failedCount := 0

	for _, result := range classifications {
		sourcePath := filepath.Join(sourceDir, result.FilePath)
		targetDir := filepath.Join(baseDir, result.Category)
		targetPath := filepath.Join(targetDir, result.FilePath)

		// 确保目标目录存在
		err := os.MkdirAll(targetDir, 0755)
		if err != nil {
			fmt.Printf("创建目标目录失败 %s: %v\n", targetDir, err)
			failedCount++
			continue
		}

		// 复制文件
		err = copyFile(sourcePath, targetPath)
		if err != nil {
			fmt.Printf("复制文件失败 %s -> %s: %v\n", sourcePath, targetPath, err)
			failedCount++
			continue
		}

		copiedCount++
	}

	fmt.Printf("\n=== 文件复制结果 ===\n")
	fmt.Printf("成功复制: %d 个文件\n", copiedCount)
	fmt.Printf("复制失败: %d 个文件\n", failedCount)
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	if err != nil {
		return err
	}

	return destFile.Sync()
}
