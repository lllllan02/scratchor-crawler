package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetView(t *testing.T) {
	t.Run("单选题", func(t *testing.T) {
		view, err := client.GetView("https://tiku.scratchor.com/question/view/abjyobpc2yhjb2ei")
		assert.NoError(t, err)
		assert.Equal(t, &View{
			Tags: []string{"2025年", "单选题", "一级", "一般"},
			Question: &Question{
				Type:  "单选",
				Alias: "abjyobpc2yhjb2ei",
				Body:  `<p>当前背景是“Arctic”，想要切换到“Bedroom 2”的背景，应使用下列哪个选项的积木？（ ）</p><p><img src="https://imgtk.scratchor.com/data/image/2025/07/14/7298_yeym_8430.png" width="87" height="297" style="width: 87px; height: 297px;"/></p>`,
				Option: []string{
					`<p><img src="https://imgtk.scratchor.com/data/image/2025/07/14/7299_fu9n_7777.png"/></p>`,
					`<p><img src="https://imgtk.scratchor.com/data/image/2025/07/14/7300_9gmh_2578.png"/></p>`,
					`<p><img src="https://imgtk.scratchor.com/data/image/2025/07/14/7301_eifg_8268.png"/></p>`,
					`<p><img src="https://imgtk.scratchor.com/data/image/2025/07/14/7302_ugv6_6854.png"/></p>`,
				},
			},
		}, view)
	})

	t.Run("多选题", func(t *testing.T) {
		view, err := client.GetView("https://tiku.scratchor.com/question/view/hrslpuuqoqvaytgk")
		assert.NoError(t, err)
		assert.Equal(t, &View{
			Tags: []string{"2025年", "多选题", "一级", "一般"},
			Question: &Question{
				Type:  "多选",
				Alias: "hrslpuuqoqvaytgk",
				Body:  `<p>下列哪2个机构可以用来省力？（ ）</p>`,
				Option: []string{
					`<p>动滑轮</p>`,
					`<p>定滑轮</p>`,
					`<p>省力杠杆</p>`,
					`<p>费力杠杆</p>`,
				},
			},
		}, view)
	})

	t.Run("判断题", func(t *testing.T) {
		view, err := client.GetView("https://tiku.scratchor.com/question/view/n70nezbpz7jhjjrc")
		assert.NoError(t, err)
		assert.Equal(t, &View{
			Tags: []string{"练习"},
			Question: &Question{
				Type:   "判断",
				Alias:  "n70nezbpz7jhjjrc",
				Body:   `<p>“脸盲症”的致病原因是患者的的眼睛有疾病。</p>`,
				Option: []string{`正确`, `错误`},
			},
		}, view)
	})

	t.Run("填空题", func(t *testing.T) {
		view, err := client.GetView("https://tiku.scratchor.com/question/view/ifufbngvurxwrtl3")
		assert.NoError(t, err)
		assert.Equal(t, &View{
			Tags: []string{"练习"},
			Question: &Question{
				Type:  "填空",
				Alias: "ifufbngvurxwrtl3",
				Body:  "<p>谁发明了蒸汽机（\u00a0 \u00a0）</p>",
			},
		}, view)
	})

	t.Run("问答题", func(t *testing.T) {
		view, err := client.GetView("https://tiku.scratchor.com/question/view/d9pvk7okovzdtu2d")
		assert.NoError(t, err)
		assert.Equal(t, &View{
			Tags: []string{"2025年", "编程题", "一级", "一般"},
			Question: &Question{
				Type:  "问答",
				Alias: "d9pvk7okovzdtu2d",
				Body:  `<p style="text-align: left;">小明去超市买了苹果和香蕉，苹果每斤6.5元，香蕉每斤4.8元。小明买了m斤苹果和n斤香蕉（m和n都是不是0的整数），请写一段程序计算小明一共需要支付多少钱？</p><p style="text-align: left;"><strong>要求：</strong></p><p style="text-align: left;">（1）程序开始运行后，需要用户输入m和n的值（整数），可以分两次输入；</p><p style="text-align: left;">（2）用户输入斤数时，要有提示语，提示语分别为：“请输入苹果斤数：”、 “请输入香蕉斤数：”；</p><p style="text-align: left;">（3）计算公式正确，正确实现总费用的计算逻辑；</p><p style="text-align: left;">（4）输出格式正确，输出字符串包含提示文本，如“小明一共需要支付：”，“元”；</p><p style="text-align: left;">（5）代码规范，运行正常。</p><p style="text-align: left;"><strong>友情提示：</strong></p><p style="text-align: left;">由于考试平台暂不支持eval()命令，同学们可以选用其他命令；当然如果您使用了，只要程序是正确的，我们阅卷时依然按照正常处理。</p>`,
			},
		}, view)
	})

	t.Run("组合题", func(t *testing.T) {
		view, err := client.GetView("https://tiku.scratchor.com/question/view/syrdjesndgxqiobb")
		assert.NoError(t, err)
		assert.Equal(t, &View{
			Tags: []string{"普及组", "初赛", "2023", "完善程序"},
			Question: &Question{
				Type:  "组合",
				Alias: "syrdjesndgxqiobb",
				Body:  "<p>（编辑距离）给定两个字符串，每次操作可以选择删除（Delete）、插入（Insert）、替换（Replace），一个字符，求将第一个字符串转换为第二个字符串所需要的最少操作次数。</p><pre class=\"brush:cpp;toolbar:false\">#include\u00a0&lt;iostream&gt;\n#include\u00a0&lt;string&gt;\n#include\u00a0&lt;vector&gt;\nusing\u00a0namespace\u00a0std;\n\nint\u00a0min(int\u00a0x,int\u00a0y,int\u00a0z){\nreturn\u00a0min(min(x,y),z);\n}\n\nint\u00a0edit_dist_dp(string\u00a0str1,string\u00a0str2){\nint\u00a0m=str1.length();\nint\u00a0n=str2.length();\nvector&lt;vector&lt;int&gt;&gt;\u00a0dp(m+1,vector&lt;int&gt;(n+1));\n\nfor(int\u00a0i=0;i&lt;=m;i++){\nfor(int\u00a0j=0;j&lt;=n;j++){\nif(i==0)\ndp[i][j]=\u00a0\u00a0\u00a0\u00a0\u00a0①\u00a0\u00a0\u00a0\u00a0\u00a0;\nelse\u00a0if(j==0)\ndp[i][j]=\u00a0\u00a0\u00a0\u00a0\u00a0②\u00a0\u00a0\u00a0\u00a0\u00a0;\nelse\u00a0if(\u00a0\u00a0\u00a0\u00a0\u00a0③\u00a0\u00a0\u00a0\u00a0\u00a0)\ndp[i][j]=\u00a0\u00a0\u00a0\u00a0\u00a0④\u00a0\u00a0\u00a0\u00a0\u00a0;\n\u00a0else\n\u00a0dp[i][j]=1+min(dp[i][j-1],dp[i-1][j],\u00a0\u00a0\u00a0\u00a0\u00a0⑤\u00a0\u00a0\u00a0\u00a0\u00a0);\u00a0\n\u00a0}\n\u00a0}\nreturn\u00a0dp[m][n];\n}\n\nint\u00a0main(){\n\u00a0string\u00a0str1,str2;\n\u00a0cin&gt;&gt;str1&gt;&gt;str2;\n\u00a0cout&lt;&lt;&#34;Mininum\u00a0number\u00a0of\u00a0operation:&#34;\n\u00a0&lt;&lt;edit_dist_dp(str1,str2)&lt;&lt;endl;\n\u00a0return\u00a00;\u00a0\n}</pre><p><br/></p>",
			},
			Items: []*Question{
				{Alias: "3tto4ramxabuto9j", Type: "单选", Body: "<p>①处应填（ ）</p>", Option: []string{
					"<p>j</p>", "<p>i</p>", "<p>m</p>", "<p>n</p>",
				}},
				{Alias: "lcdncixsa8tarmxb", Type: "单选", Body: "<p>②处应填（ ）</p>", Option: []string{
					"<p>j</p>", "<p>i</p>", "<p>m</p>", "<p>n</p>",
				}},
				{Alias: "czso7sqyx5vzjebs", Type: "单选", Body: "<p>③处应填（ ）</p>", Option: []string{
					"<p>str1[i-1]==str2[j-1]</p>",
					"<p>str1[i]==str2[j]</p>",
					"<p>str1[i-1]!=str2[j-1]</p>",
					"<p>str1[i]!=str2[j]</p>",
				}},
				{Alias: "dnybf4phfbrvfcm3", Type: "单选", Body: "<p>④处应填（ ）</p>", Option: []string{
					"<p>dp[i-1][j-1]+1</p>",
					"<p>dp[i-1][j-1]</p>",
					"<p>dp[i-1][j]</p>",
					"<p>dp[i][j-1]</p>",
				}},
				{Alias: "r8vrtlegrfmmpw8v", Type: "单选", Body: "<p>⑤处应填（ ）</p>", Option: []string{
					"<p>dp[i][j] + 1</p>",
					"<p>dp[i-1][j-1]+1</p>",
					"<p>dp[i-1][j-1]</p>",
					"<p>dp[i][j]</p>",
				}},
			},
		}, view)
	})
}
