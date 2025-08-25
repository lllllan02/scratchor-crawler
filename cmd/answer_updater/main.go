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
	if client, err = api.NewClient("Hm_lvt_22880f1d42e3788b94e7aa1361ea923a=1755757549; HMACCOUNT=4D5E77D5BB397FC9; 333decb516a63de949aa73f356ff0515=51967059785189c5bd710720adcf8afb; ssid=eyJpdiI6ImFLeDVPd3AyVVExVlR2VzhpRGVTdXc9PSIsInZhbHVlIjoiNXF3Wk10TE40OUl2OEhCY252cFVjT0dZVmd2eUdrRlNiXC85QjY4Q3Ztdk5HWTFLaTlVTjJSS1hZcUxrd1RDNXZSZ3BqOVN2XC8yc1NSTFBtdlBBTjdpdz09IiwibWFjIjoiMDU3MzllNjBmMmZlYzFjZTQ5N2EwMmI1NjdkNmE3N2RkNTY5ZDM2YjAyOGY2ZjhlZDg2YWM3MzkzODdiMjMwMiJ9; Hm_lpvt_22880f1d42e3788b94e7aa1361ea923a=1755854532"); err != nil {
		fmt.Printf("%s创建客户端失败%s: %v\n", utils.ColorRed, utils.ColorReset, err)
		os.Exit(1)
	}

	utils.ProcessFiles("data", func(view *api.View) (needSave bool, err error) {
		for _, item := range view.Items {
			ns, err := updateAnswer(item)
			if err != nil {
				return false, err
			}

			needSave = needSave || ns
		}

		ns, err := updateAnswer(view.Question)
		if err != nil {
			return false, err
		}
		return needSave || ns, nil
	})
}

func updateAnswer(question *api.Question) (bool, error) {
	if question.Answer != nil || question.Analysis != "" {
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
