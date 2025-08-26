package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/lllllan02/scratchor-crawler/api"
	"github.com/lllllan02/scratchor-crawler/utils"
)

// 统计结果结构
type StatResult struct {
	TotalQuestions  int            `json:"total_questions"`
	AnsweredCount   int            `json:"answered_count"`
	UnansweredCount int            `json:"unanswered_count"`
	QuestionTypes   map[string]int `json:"question_types"`
}

func main() {
	fmt.Println("开始统计题目数据...")

	var totalQuestions, totalAnswered int
	questionTypes := make(map[string]int)

	// 处理所有文件
	utils.ProcessFiles("data", func(view *api.View) (needSave bool, err error) {
		// 统计主题目
		totalQuestions++

		// 统计题目类型
		if view.Type != "" {
			questionTypes[view.Type]++
		}

		// 检查是否有答案
		if hasAnswer(view.Question) {
			totalAnswered++
		}

		return false, nil
	})

	// 创建统计结果
	result := StatResult{
		TotalQuestions:  totalQuestions,
		AnsweredCount:   totalAnswered,
		UnansweredCount: totalQuestions - totalAnswered,
		QuestionTypes:   questionTypes,
	}

	// 将结果写入文件
	outputPath := "statistics.json"
	file, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("创建文件失败: %v\n", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(result); err != nil {
		fmt.Printf("写入JSON失败: %v\n", err)
		return
	}

	fmt.Printf("统计完成！结果已保存到: %s\n", outputPath)
	fmt.Printf("总题目数: %d\n", totalQuestions)
	fmt.Printf("已答题目: %d\n", totalAnswered)
	fmt.Printf("未答题目: %d\n", totalQuestions-totalAnswered)
	fmt.Printf("答题率: %.1f%%\n", float64(totalAnswered)/float64(totalQuestions)*100)
}

// 检查题目是否有答案
func hasAnswer(question *api.Question) bool {
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
