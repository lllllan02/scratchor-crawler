package main

import (
	"flag"
	"fmt"

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
	if client, err = api.NewClient("HMACCOUNT=4D5E77D5BB397FC9; Hm_lvt_22880f1d42e3788b94e7aa1361ea923a=1758358879; 7c5ea60f97881a7fc74285aca53772d3=9045a006eb06c91999484624d03ba099; ssid=eyJpdiI6IndaWmFLenhwalhkOEdTVGpaOVNSRUE9PSIsInZhbHVlIjoid0lrZEdIMEkwdmxxWmQ0RG5EN3JSSHpLY0hxUWVOZkxtcVwvZ2tuNnQxdTN0SXZKQm93TzRMK1RZVm1ZQVhqMlwvWVlRMWNEREtpRjlKWWhvb21pZ0Vjdz09IiwibWFjIjoiNTdiZTZkOTY1MDU3NmU1NzhkYmU1YmNlYmRmYTBkOWIyMjNkMzg3NjZkYzVjODlmNjYyNWM0NDA4YTY3MWZmNCJ9; Hm_lpvt_22880f1d42e3788b94e7aa1361ea923a=1758675651"); err != nil {
		panic(err)
	}

	Categories()

	utils.WriteJSON("paper_map.json", categories)
}

type Category struct {
	Alias    string
	Url      string
	Title    string
	Children []*Category
}

var categories = []*Category{
	{
		"",
		"",
		"电子学会考级",
		[]*Category{
			{
				"",
				"",
				"Scratch等级考试",
				[]*Category{
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=1",
						"Scratch一级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=2",
						"Scratch二级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=3",
						"Scratch三级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=4",
						"Scratch四级模拟",
						nil,
					},
				},
			},
			{
				"",
				"",
				"Python等级考试",
				[]*Category{
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=5",
						"Python一级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=6",
						"Python二级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=7",
						"Python三级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=8",
						"Python四级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=9",
						"Python五级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=19",
						"Python六级模拟",
						nil,
					},
				},
			},
			{
				"",
				"",
				"C/C++等级考试",
				[]*Category{
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=24",
						"C/C++一级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=25",
						"C/C++二级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=26",
						"C/C++三级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=27",
						"C/C++四级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=28",
						"C/C++五级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=29",
						"C/C++六级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=30",
						"C/C++七级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=31",
						"C/C++八级模拟",
						nil,
					},
				},
			},
			{
				"",
				"",
				"机器人等级考试",
				[]*Category{
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=10",
						"机器人一级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=11",
						"机器人二级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=12",
						"机器人三级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=13",
						"机器人四级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=14",
						"机器人五级模拟",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=15",
						"机器人六级模拟",
						nil,
					},
				},
			},
		},
	},
	{
		"",
		"",
		"蓝桥竞赛",
		[]*Category{
			{
				"",
				"",
				"蓝桥Scratch",
				[]*Category{
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=55",
						"STEMA",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=56",
						"省赛",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=57",
						"国赛",
						nil,
					},
				},
			},
			{
				"",
				"",
				"蓝桥Python",
				[]*Category{
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=58",
						"STEMA",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=59",
						"省赛",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=60",
						"国赛",
						nil,
					},
				},
			},
			{
				"",
				"",
				"蓝桥C++",
				[]*Category{
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=61",
						"STEMA",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=62",
						"省赛",
						nil,
					},
					{
						"",
						"https://tiku.scratchor.com/paper?categoryId=63",
						"国赛",
						nil,
					},
				},
			},
			{
				"",
				"https://tiku.scratchor.com/paper?categoryId=35",
				"科技素养",
				nil,
			},
			{
				"",
				"https://tiku.scratchor.com/paper?categoryId=36",
				"计算思维",
				nil,
			},
		},
	},
	{
		"",
		"",
		"信息素养大赛",
		[]*Category{
			{
				"",
				"https://tiku.scratchor.com/paper?categoryId=41",
				"图形化编程挑战赛",
				nil,
			},
			{
				"",
				"https://tiku.scratchor.com/paper?categoryId=42",
				"Python编程挑战赛",
				nil,
			},
			{
				"",
				"https://tiku.scratchor.com/paper?categoryId=52",
				"C++编程挑战赛",
				nil,
			},
		},
	},
	{
		"",
		"https://tiku.scratchor.com/paper?categoryId=17",
		"信息学奥赛",
		nil,
	},
	{
		"",
		"",
		"CCF编程能力等级认证（GESP）",
		[]*Category{
			{
				"",
				"https://tiku.scratchor.com/paper?categoryId=44",
				"图形化编程",
				nil,
			},
			{
				"",
				"https://tiku.scratchor.com/paper?categoryId=45",
				"Python编程",
				nil,
			},
			{
				"",
				"https://tiku.scratchor.com/paper?categoryId=46",
				"C++编程",
				nil,
			},
		},
	},
}

func Categories() error {
	for _, category := range categories {
		if err := GetCategory(category); err != nil {
			return err
		}
	}
	return nil
}

func GetCategory(category *Category) error {
	for _, child := range category.Children {
		if err := GetCategory(child); err != nil {
			return err
		}
	}

	var links []string
	var err error
	for next := category.Url; next != ""; {
		links, next, err = client.GetPaperList(next)
		if err != nil {
			return err
		}

		for _, link := range links {
			var alias string
			fmt.Sscanf(link, "https://tiku.scratchor.com/paper/view/%s", &alias)
			category.Children = append(category.Children, &Category{Alias: alias})
		}
	}

	return nil
}
