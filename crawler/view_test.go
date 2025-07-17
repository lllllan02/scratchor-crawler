package crawler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetView(t *testing.T) {
	for _, url := range []string{
		"https://tiku.scratchor.com/question/view/abjyobpc2yhjb2ei", // 单选题
		"https://tiku.scratchor.com/question/view/hrslpuuqoqvaytgk", // 多选题
		"https://tiku.scratchor.com/question/view/n70nezbpz7jhjjrc", // 判断题
		"https://tiku.scratchor.com/question/view/ifufbngvurxwrtl3", // 填空题
		"https://tiku.scratchor.com/question/view/d9pvk7okovzdtu2d", // 问答题
		"https://tiku.scratchor.com/question/view/syrdjesndgxqiobb", // 组合题
	} {
		view, err := GetView(url, cookie)
		assert.NoError(t, err)

		fmt.Printf("view.Tags: %+v\n", view.Tags)
		fmt.Printf("view.Question: %+v\n", view.Question)
		for _, item := range view.Items {
			fmt.Printf("item: %+v\n", item)
		}
	}
}
