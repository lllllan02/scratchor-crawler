package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

// WriteJSON 将对象序列化为JSON并写入文件
// filePath: 目标文件路径
// obj: 要序列化的对象
func WriteJSON(filePath string, obj any) error {
	// 将对象转换为JSON，不对字符进行转义
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")

	if err := encoder.Encode(obj); err != nil {
		return fmt.Errorf("JSON序列化失败: %v", err)
	}

	data := buf.Bytes()

	// 写入文件
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// ReadJSON 从文件读取JSON内容到指定类型的对象中
// filePath: 文件路径
// 返回指定类型的对象和错误
func ReadJSON[T any](filePath string) (T, error) {
	var obj T

	// 读取文件内容
	data, err := os.ReadFile(filePath)
	if err != nil {
		return obj, fmt.Errorf("读取文件失败: %v", err)
	}

	// 解析JSON到对象
	if err := json.Unmarshal(data, &obj); err != nil {
		return obj, fmt.Errorf("JSON解析失败: %v", err)
	}

	return obj, nil
}
