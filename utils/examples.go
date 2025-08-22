package utils

import (
	"fmt"
	"strings"
)

// 示例：创建一个简单的文本替换处理器
func CreateTextReplacer(oldText, newText string) func(string) (bool, error) {
	return func(filePath string) (bool, error) {
		// 读取文件内容
		content, err := ReadFile(filePath)
		if err != nil {
			return false, fmt.Errorf("读取文件失败 %s: %w", filePath, err)
		}

		// 检查是否需要替换
		contentStr := string(content)
		if !strings.Contains(contentStr, oldText) {
			return false, nil
		}

		// 执行替换
		newContentStr := strings.ReplaceAll(contentStr, oldText, newText)
		newContent := []byte(newContentStr)

		// 保存文件
		err = WriteFile(filePath, newContent)
		if err != nil {
			return false, fmt.Errorf("保存文件失败 %s: %w", filePath, err)
		}

		return true, nil
	}
}

// 示例：创建一个文件统计处理器
func CreateFileAnalyzer() func(string) (bool, error) {
	return func(filePath string) (bool, error) {
		// 读取文件内容
		content, err := ReadFile(filePath)
		if err != nil {
			return false, fmt.Errorf("读取文件失败 %s: %w", filePath, err)
		}

		// 分析文件内容（这里只是示例，不实际修改文件）
		contentStr := string(content)
		lineCount := strings.Count(contentStr, "\n") + 1
		charCount := len(contentStr)
		wordCount := len(strings.Fields(contentStr))

		// 输出统计信息（在实际使用中，你可能想要记录到日志文件）
		fmt.Printf("文件: %s, 行数: %d, 字符数: %d, 单词数: %d\n",
			filePath, lineCount, charCount, wordCount)

		// 不修改文件，返回 false
		return false, nil
	}
}

// 示例：创建一个条件处理器
// 这个处理器只在满足特定条件时才处理文件
func CreateConditionalProcessor(condition func(string) bool, processor func(string) (bool, error)) func(string) (bool, error) {
	return func(filePath string) (bool, error) {
		// 检查条件
		if !condition(filePath) {
			return false, nil
		}

		// 执行实际处理
		return processor(filePath)
	}
}

// 示例：创建一个批量处理器
// 这个处理器可以同时执行多个处理函数
func CreateBatchProcessor(processors ...func(string) (bool, error)) func(string) (bool, error) {
	return func(filePath string) (bool, error) {
		needSave := false

		// 依次执行所有处理器
		for _, processor := range processors {
			save, err := processor(filePath)
			if err != nil {
				return false, err
			}
			if save {
				needSave = true
			}
		}

		return needSave, nil
	}
}
