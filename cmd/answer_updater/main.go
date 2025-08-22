package main

import (
	"fmt"
	"os"

	"github.com/lllllan02/scratchor-crawler/api"
	"github.com/lllllan02/scratchor-crawler/utils"
)

var client *api.Client

func main() {
	var err error
	if client, err = api.NewClient("9c8adef9772544421fc21404af7ed346=e6f5696bd5ff493578540e8afee89af8; ssid=eyJpdiI6IjBRYUlUYzdNeFM4RGhjZ2pOT0RNWnc9PSIsInZhbHVlIjoiZUUzSW5NcTU0bFwvaVYwRFwvaVFIRXhzZldUYWFtcVdVMkpETlZKQmxTb1lkaWNKSGxmT05lVmdNMjA1XC8rQVZhXC84SjY1TG5aWGJRSGlkU0lUblJaZjhBPT0iLCJtYWMiOiI3YWVmYTk0NDQxZjg3YTFiMDlkYjhlNWUwMDNmMzFhM2FkNjBmZWY0Yzg1OWVhYjQ1OTQ0M2E3ODRhNjkzMzI0In0%3D"); err != nil {
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
