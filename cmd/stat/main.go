package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/lllllan02/scratchor-crawler/api"
	"github.com/lllllan02/scratchor-crawler/utils"
)

// 目录统计结构
type DirectoryNode struct {
	Name            string                    `json:"name"`
	Path            string                    `json:"path"`
	TotalQuestions  int                       `json:"total_questions"`
	AnsweredCount   int                       `json:"answered_count"`
	UnansweredCount int                       `json:"unanswered_count"`
	AnswerRate      float64                   `json:"answer_rate"`
	QuestionTypes   map[string]int            `json:"question_types"`
	SubDirectories  map[string]*DirectoryNode `json:"sub_directories,omitempty"`
}

func main() {
	fmt.Println("开始统计题目数据...")

	// 初始化根目录统计
	rootStats := &DirectoryNode{
		Name:           "data",
		Path:           "data",
		QuestionTypes:  make(map[string]int),
		SubDirectories: make(map[string]*DirectoryNode),
	}

	// 处理所有文件
	utils.ProcessFiles("data", func(filePath string, view *api.Question) (needSave bool, err error) {
		// 更新根目录统计
		rootStats.TotalQuestions++
		if hasAnswer(view.QuestionBody) {
			rootStats.AnsweredCount++
		}
		if view.Type != "" {
			rootStats.QuestionTypes[view.Type]++
		}

		// 获取文件所在目录路径并更新目录树
		dirPath := getDirectoryFromPath(filePath)
		if dirPath != "" {
			updateDirectoryTree(rootStats, dirPath, view.Type, hasAnswer(view.QuestionBody))
		}

		return false, nil
	})

	// 计算所有目录的答题率和未答题数
	calculateDirectoryStats(rootStats)

	// 保存结果
	utils.WriteJSON("stat.json", rootStats)

	// 打印统计摘要
	printDirectoryStats(rootStats)
}

// 更新目录树统计
func updateDirectoryTree(root *DirectoryNode, dirPath, questionType string, hasAnswer bool) {
	parts := strings.Split(dirPath, "/")
	current := root

	// 逐级创建或更新目录节点
	for i, part := range parts {
		if part == "" {
			continue
		}

		// 确保子目录存在
		if current.SubDirectories[part] == nil {
			current.SubDirectories[part] = &DirectoryNode{
				Name:           part,
				Path:           strings.Join(parts[:i+1], "/"),
				QuestionTypes:  make(map[string]int),
				SubDirectories: make(map[string]*DirectoryNode),
			}
		}

		// 更新当前目录统计
		dir := current.SubDirectories[part]
		dir.TotalQuestions++
		if hasAnswer {
			dir.AnsweredCount++
		}
		if questionType != "" {
			dir.QuestionTypes[questionType]++
		}

		current = dir
	}
}

// 计算目录统计数据的答题率和未答题数
func calculateDirectoryStats(node *DirectoryNode) {
	node.UnansweredCount = node.TotalQuestions - node.AnsweredCount
	if node.TotalQuestions > 0 {
		node.AnswerRate = float64(node.AnsweredCount) / float64(node.TotalQuestions) * 100
	}

	// 递归计算子目录
	for _, subDir := range node.SubDirectories {
		calculateDirectoryStats(subDir)
	}
}

// 从文件路径获取目录路径
func getDirectoryFromPath(filePath string) string {
	relPath, err := filepath.Rel("data", filePath)
	if err != nil {
		return ""
	}
	return filepath.Dir(relPath)
}

// 打印目录统计信息
func printDirectoryStats(node *DirectoryNode) {
	printNodeStats(node, 0)
	fmt.Printf("\n统计完成！结果已保存到: stat.json\n")
}

// 递归打印节点统计信息
func printNodeStats(node *DirectoryNode, level int) {
	indent := strings.Repeat("  ", level)

	if level == 0 {
		fmt.Printf("\n=== 整体统计 ===\n")
	} else {
		fmt.Printf("%s📁 %s: %d题 (答题率: %.1f%%)\n",
			indent, node.Name, node.TotalQuestions, node.AnswerRate)
	}

	// 打印题目类型分布
	if len(node.QuestionTypes) > 0 {
		fmt.Printf("%s  📊 题目类型: ", indent)
		for questionType, count := range node.QuestionTypes {
			fmt.Printf("%s(%d) ", questionType, count)
		}
		fmt.Println()
	}

	// 递归打印子目录
	for _, subDir := range node.SubDirectories {
		printNodeStats(subDir, level+1)
	}
}

// 检查题目是否有答案
func hasAnswer(question *api.QuestionBody) bool {
	// 检查Answer字段
	if len(question.Answer) > 0 {
		for _, answer := range question.Answer {
			if strings.TrimSpace(answer) != "" {
				return true
			}
		}
	}

	// 检查Analysis字段
	if strings.TrimSpace(question.Analysis) != "" {
		return true
	}

	return false
}
