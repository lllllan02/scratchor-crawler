package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAnswer(t *testing.T) {
	t.Run("单选题", func(t *testing.T) {
		result, err := client.GetAnswer("abjyobpc2yhjb2ei", "单选题")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []string{"D"}, result.Answer)
		assert.Equal(t, "<p>确定要切换到“Bedroom 2”背景，选择换成“Bedroom 2”背景即可，其他操作都不能切换到“Bedroom 2”背景。</p>", result.Analysis)
	})

	t.Run("判断题", func(t *testing.T) {
		result, err := client.GetAnswer("btbzgxyhcetfpiyf", "判断题")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []string{"A"}, result.Answer)
		assert.Equal(t, "<p>当角色只有一个造型时，该造型没有删除按钮。</p>", result.Analysis)
	})

	t.Run("多选题", func(t *testing.T) {
		result, err := client.GetAnswer("n7ycxzyef8uxsyts", "多选题")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []string{"A", "B", "C"}, result.Answer)
		assert.Equal(t, "", result.Analysis)
	})

	t.Run("填空题", func(t *testing.T) {
		result, err := client.GetAnswer("sahgxbamvtvpbsnn", "填空题")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []string{"去了", "没去", "没去", "没下雨"}, result.Answer)
		assert.Equal(t, "", result.Analysis)
	})

	t.Run("问答题", func(t *testing.T) {
		result, err := client.GetAnswer("7w3z8vedemkz2nmf", "问答题")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, []string{`<p style="text-align: left;"><strong>答案解析：</strong></p><p style="text-align: left;"><strong>Glow-1</strong></p><p><img src="https://imgtk.scratchor.com/data/image/2025/07/14/3308_me1e_7340.png" width="164" height="300" style="width: 164px; height: 300px;"/></p><p style="text-align: left;"><strong>背景</strong></p><p><img src="https://imgtk.scratchor.com/data/image/2025/07/14/3309_tzh2_1181.png" width="177" height="214" style="width: 177px; height: 214px;"/></p>`}, result.Answer)
		assert.Equal(t, `<p style="text-align: left;"><strong>评分标准:</strong></p><p>（1）正确添加两个背景；（2分）<br/></p><p>（2）正确为背景添加声音：Bossa Nova；（1分）<br/></p><p>（3）删除小猫角色；（1分）<br/></p><p>（4）正确添加角色Glow-1（1分），正确添加造型Glow-2、Glow-3；（2分）<br/></p><p style="text-align: left;">（5）角色初始位置在舞台中心，初始造型为Glow-3；（2分）</p><p>（6）初始背景为Hearts；（1分）<br/></p><p>（7）等待1秒，角色造型切换为Glow-2；（1分）<br/></p><p>（8）等待1秒，角色造型切换为Glow-1，背景切换为Theater；（2分）<br/></p><p>（9）声音播放正确。（2分）</p>`, result.Analysis)
	})

	t.Run("组合题", func(t *testing.T) {
		result, err := client.GetAnswer("syrdjesndgxqiobb", "组合题")
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Nil(t, result.Answer)
		assert.Equal(t, `<p>先阅读 main 函数，可以得知代码读入了两个字符串之后输出了 edit_dist_dp 的返回值，因此重点关注 edit_dist_dp 函数；</p><p>11-13 行定义了 n 为 str1 的长度，m 为 str2 的长度，一个二维的 vector（可以当成二维数组）dp[0..m][0..n]；</p><p>15-26 行为 dp 过程；</p><p>27 行为返回 dp[n][m]，看完这一行就应该明白 dp 数组的含义，之后通过 dp 状态将 15-26 行的空填出来；</p><p>dp[n][m] 表示将 str1 长度为 n 的前缀变成 str2 长度为 m 的前缀需要的最少步数，特别注意 str1, str2 是 string，下标范围分别是 [0..m-1] 与 [0..n-1]。</p><p>（1）考察 dp[0][j] 的值，根据阅读得到的 dp 数组的含义，dp[0][j] 表示 str1 长度为 0 的前缀（空字符串）变成 str2 长度为 j 的前缀的最少步数，这里明显最少在空串内添加 str2 的前 j 个字符，因此填 j选 A 选项</p><p>（2）考察 dp[i][0] 的值，根据阅读得到的 dp 数组的含义，dp[i][0] 表示 str1 长度为 i 的前缀变成str2 长度为 0 的前缀（空字符串）的最少步数，这里明显最少是将 str1 的前 i个字符全部删除，因此填 i 选 B 选项</p><p>（3）考察编辑距离 dp 转移，需要阅读 22-24 行的代码可知这里应该是两个字符相等不需要操作的情况，也就是 str1 的第 i个字符与 str2 的第 j 个字符相等的情况，但是要特别注意的是这一问埋了大坑，str1 的第 i个字符是str1[i-1]，不要跳进去了，选 A</p><p>（4）40题时已经说过，如果两个字符相等的话，不需要操作，因此将 str1 前 i个字符变成 str2 前 j 个字符需要的最少步数就与将 str1 前 i-1 个字符变成 str2 前j-1 个字符是一样的，填 dp[i-1][j-1]，选 B</p><p>（5）这里有一个在上面定义的 min 函数，功能是对三个整数求最小值，观察第 24 行 dp[i][j] = 1 + min(dp[i][j - 1], dp[i - 1][j], ?) 由前面的 1 + 可知这里进行了一次操作，那么dp[i][j - 1] 就对应着插入，dp[i - 1][j] 对应着删除，剩下要填的就是替换了，填 dp[i-1][j-1]，选 C</p>`, result.Analysis)
	})
}
