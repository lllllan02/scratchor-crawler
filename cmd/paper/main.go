package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/lllllan02/scratchor-crawler/api"
	"github.com/lllllan02/scratchor-crawler/utils"
	"github.com/schollz/progressbar/v3"
)

var client *api.Client

func main() {
	// 解析命令行参数
	var limiterConcurrency uint
	flag.UintVar(&limiterConcurrency, "limiter", 2, "限流器并发数，默认为2")
	flag.Parse()

	// 设置全局限流器并发数
	api.LimiterConcurrency = limiterConcurrency

	var err error
	if client, err = api.NewClient("HMACCOUNT=4D5E77D5BB397FC9; Hm_lvt_22880f1d42e3788b94e7aa1361ea923a=1758358879; a0e0cf20f929d48c545b6794f503b8d2=14ea9ecc4ecea3c94148df41959a843d; ssid=eyJpdiI6ImdxUGpWRkhGMWdyZmwzXC9XSElhKzlRPT0iLCJ2YWx1ZSI6ImNqemx6Y2poNFNMU2tLWE1wc0VhSFVMSFVXOHM2S090bUtLK3JJbTI4bjFyYnpJMk5BXC9vUzZYXC9MVU9pRDNaNGIralJYMDAyU1h3OVpWa1VJOEdBY1E9PSIsIm1hYyI6ImUzZDk0NDU1NzUyZmJjYjU1N2EyYzEzMzA0NmM5ZGJkYzFkMGYwYTA2NTFhNWEzMzI4MmY2OGI2NmFkZTI1MmEifQ%3D%3D; Hm_lpvt_22880f1d42e3788b94e7aa1361ea923a=1758502503"); err != nil {
		panic(err)
	}

	if err := Papers(); err != nil {
		fmt.Printf("%s程序执行失败: %v%s\n", utils.ColorRed, err, utils.ColorReset)
		os.Exit(1)
	}
}

func Papers() error {
	fmt.Printf("%s开始下载试卷...%s\n", utils.ColorCyan, utils.ColorReset)

	// 创建 paper 目录
	os.MkdirAll("paper", 0755)

	// 创建进度条（使用 -1 表示未知总数）
	bar := progressbar.NewOptions(
		1271,
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(50),
		progressbar.OptionSetDescription(fmt.Sprintf("%s下载试卷%s", utils.ColorCyan, utils.ColorReset)),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	// 直接开始下载，边获取边下载
	var downloadedCount int
	for next := "https://tiku.scratchor.com/paper"; next != ""; {
		pageLinks, nextPage, err := client.GetPaperList(next)
		if err != nil {
			fmt.Printf("%s获取试卷列表失败: %v%s\n", utils.ColorRed, err, utils.ColorReset)
			return err
		}

		// 下载当前页面的试卷
		for _, link := range pageLinks {
			if err := Paper(link, bar, downloadedCount+1, -1); err != nil {
				return err
			}
			downloadedCount++
			bar.Add(1)
		}

		next = nextPage
	}

	bar.Finish()
	fmt.Printf("\n%s下载完成，共下载 %d 个试卷%s\n", utils.ColorGreen, downloadedCount, utils.ColorReset)
	return nil
}

func Paper(url string, bar *progressbar.ProgressBar, current, total int) error {
	// 从 URL 中提取 paper ID
	var paperID string
	if _, err := fmt.Sscanf(url, "https://tiku.scratchor.com/paper/view/%s", &paperID); err != nil {
		fmt.Printf("%s解析 URL 中的 paper ID 失败: %v%s\n", utils.ColorRed, err, utils.ColorReset)
		return err
	}

	// 构建文件路径
	filePath := fmt.Sprintf("paper/%s.json", paperID)

	// 检查文件是否已存在
	if _, err := os.Stat(filePath); err == nil {
		progressbar.Bprintln(bar, fmt.Sprintf("试卷「%s」已存在，跳过", paperID))
		return nil
	}

	// 获取试卷信息
	paper, err := client.GetPaper(url)
	if err != nil {
		fmt.Printf("%s获取试卷「%s」失败: %v%s\n", utils.ColorRed, url, err, utils.ColorReset)
		return err
	}

	if paper.Head == "" {
		fmt.Printf("%s获取试卷「%s」失败: %v%s\n", utils.ColorRed, url, err, utils.ColorReset)
		return errors.New("获取试卷信息失败")
	}

	// 显示试卷信息
	var progressInfo string
	if total > 0 {
		progressInfo = fmt.Sprintf("[%d/%d] ", current, total)
	}
	progressbar.Bprintln(bar, fmt.Sprintf("%s获取试卷「%s」: %s (%s题, %s分, %s)",
		progressInfo, paperID, paper.Head, paper.TotalQuestions, paper.TotalScore, paper.Duration))

	// 将试卷信息保存为 JSON
	if err := utils.WriteJSON(filePath, paper); err != nil {
		fmt.Printf("%s保存试卷「%s」失败: %v%s\n", utils.ColorRed, paperID, err, utils.ColorReset)
		return err
	}

	return nil
}
