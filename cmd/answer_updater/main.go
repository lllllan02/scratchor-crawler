package main

import (
	"fmt"
	"os"

	"github.com/lllllan02/scratchor-crawler/api"
	"github.com/lllllan02/scratchor-crawler/utils"
)

func main() {
	// 创建 API 客户端
	client, err := api.NewClient("9c8adef9772544421fc21404af7ed346=e6f5696bd5ff493578540e8afee89af8; ssid=eyJpdiI6IjBRYUlUYzdNeFM4RGhjZ2pOT0RNWnc9PSIsInZhbHVlIjoiZUUzSW5NcTU0bFwvaVYwRFwvaVFIRXhzZldUYWFtcVdVMkpETlZKQmxTb1lkaWNKSGxmT05lVmdNMjA1XC8rQVZhXC84SjY1TG5aWGJRSGlkU0lUblJaZjhBPT0iLCJtYWMiOiI3YWVmYTk0NDQxZjg3YTFiMDlkYjhlNWUwMDNmMzFhM2FkNjBmZWY0Yzg1OWVhYjQ1OTQ0M2E3ODRhNjkzMzI0In0%3D")
	if err != nil {
		fmt.Printf("%s创建客户端失败%s: %v\n", utils.ColorRed, utils.ColorReset, err)
		os.Exit(1)
	}

	// 配置答案更新器
	config := utils.AnswerUpdaterConfig{
		DataDir: "data",
		Client:  client,
	}

	// 创建答案更新器处理函数
	answerHandler := utils.CreateAnswerUpdater(config)

	// 使用通用文件遍历器处理文件
	fmt.Printf("%s开始答案更新任务%s\n", utils.ColorCyan, utils.ColorReset)
	err = utils.ProcessFiles("data", answerHandler)
	if err != nil {
		fmt.Printf("%s答案更新失败%s: %v\n", utils.ColorRed, utils.ColorReset, err)
		os.Exit(1)
	}

	fmt.Printf("%s答案更新任务完成%s\n", utils.ColorGreen, utils.ColorReset)
}
