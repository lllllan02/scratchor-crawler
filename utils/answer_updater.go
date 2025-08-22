package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lllllan02/scratchor-crawler/api"
)

// 答案更新器配置
type AnswerUpdaterConfig struct {
	DataDir string      // 数据目录
	Client  *api.Client // API客户端
}

// 创建答案更新器处理函数
func CreateAnswerUpdater(config AnswerUpdaterConfig) func(string) (bool, error) {
	return func(filePath string) (bool, error) {
		// 读取文件内容
		content, err := ReadFile(filePath)
		if err != nil {
			return false, fmt.Errorf("读取文件失败 %s: %w", filePath, err)
		}

		// 尝试解析 JSON
		var jsonData map[string]interface{}
		if err := json.Unmarshal(content, &jsonData); err != nil {
			// 如果不是 JSON 文件，跳过
			return false, nil
		}

		// 检查是否需要更新答案
		needUpdate := false

		// 示例：检查是否有特定的字段需要更新
		// 这里可以根据实际需求修改逻辑
		if question, exists := jsonData["question"]; exists {
			questionStr, ok := question.(string)
			if ok && strings.Contains(questionStr, "需要答案") {
				// 调用 API 获取答案
				answer, err := getAnswerFromAPI(config.Client, questionStr)
				if err != nil {
					return false, fmt.Errorf("获取答案失败: %w", err)
				}

				// 更新答案字段
				jsonData["answer"] = answer
				needUpdate = true
			}
		}

		// 如果有更新，保存文件
		if needUpdate {
			var buf bytes.Buffer
			encoder := json.NewEncoder(&buf)
			encoder.SetEscapeHTML(false)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(jsonData); err != nil {
				return false, fmt.Errorf("序列化 JSON 失败 %s: %w", filePath, err)
			}

			err = WriteFile(filePath, buf.Bytes())
			if err != nil {
				return false, fmt.Errorf("更新文件失败 %s: %w", filePath, err)
			}
		}

		return needUpdate, nil
	}
}

// 从 API 获取答案（示例函数）
func getAnswerFromAPI(client *api.Client, question string) (string, error) {
	// 这里实现具体的 API 调用逻辑
	// 例如：调用 AI 服务获取答案
	// 现在返回一个示例答案
	return "这是一个示例答案", nil
}

// 创建通用的内容更新器
// 这个函数接受一个自定义的更新逻辑
func CreateContentUpdater(updateFunc func(map[string]interface{}) (bool, error)) func(string) (bool, error) {
	return func(filePath string) (bool, error) {
		// 读取文件内容
		content, err := ReadFile(filePath)
		if err != nil {
			return false, fmt.Errorf("读取文件失败 %s: %w", filePath, err)
		}

		// 尝试解析 JSON
		var jsonData map[string]interface{}
		if err := json.Unmarshal(content, &jsonData); err != nil {
			// 如果不是 JSON 文件，跳过
			return false, nil
		}

		// 调用自定义更新逻辑
		needUpdate, err := updateFunc(jsonData)
		if err != nil {
			return false, fmt.Errorf("更新内容失败: %w", err)
		}

		// 如果有更新，保存文件
		if needUpdate {
			var buf bytes.Buffer
			encoder := json.NewEncoder(&buf)
			encoder.SetEscapeHTML(false)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(jsonData); err != nil {
				return false, fmt.Errorf("序列化 JSON 失败 %s: %w", filePath, err)
			}

			err = WriteFile(filePath, buf.Bytes())
			if err != nil {
				return false, fmt.Errorf("更新文件失败 %s: %w", filePath, err)
			}
		}

		return needUpdate, nil
	}
}
