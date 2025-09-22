package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/lllllan02/scratchor-crawler/api"
	"github.com/lllllan02/scratchor-crawler/utils"
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
		fmt.Printf("%s创建客户端失败%s: %v\n", utils.ColorRed, utils.ColorReset, err)
		os.Exit(1)
	}

	utils.ProcessFiles("data", func(_ string, view *api.Question) (needSave bool, err error) {
		for _, item := range view.Items {
			ns, err := updateAnswer(item)
			if err != nil {
				return false, err
			}

			needSave = needSave || ns
		}

		ns, err := updateAnswer(view.QuestionBody)
		if err != nil {
			return false, err
		}
		return needSave || ns, nil
	})
}

func updateAnswer(question *api.QuestionBody) (bool, error) {
	if question.Alias == "" || question.Answer != nil || question.Analysis != "" {
		return false, nil
	}

	answer, err := client.GetAnswer(question.Alias, question.Type)
	if err != nil {
		fmt.Printf("%s获取题目答案「%s」失败%s: %v\n", utils.ColorRed, question.Alias, utils.ColorReset, err)
		return false, err
	}

	question.Answer = answer.Answer
	question.Analysis = answer.Analysis

	return true, nil
}
