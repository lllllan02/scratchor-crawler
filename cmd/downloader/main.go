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

	// 配置图片下载器
	config := utils.ImageDownloaderConfig{
		DataDir: "data",
		ImgDir:  "image",
		Client:  client,
	}

	// 创建图片下载器处理函数
	imageHandler := utils.CreateImageDownloader(config)

	// 使用通用文件遍历器处理文件
	fmt.Printf("%s开始图片下载任务%s\n", utils.ColorCyan, utils.ColorReset)
	err = utils.ProcessFiles("data", imageHandler)
	if err != nil {
		fmt.Printf("%s图片下载失败%s: %v\n", utils.ColorRed, utils.ColorReset, err)
		os.Exit(1)
	}

	fmt.Printf("%s图片下载任务完成%s\n", utils.ColorGreen, utils.ColorReset)
}
