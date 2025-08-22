package utils

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/schollz/progressbar/v3"
)

// 定义颜色
const (
	ColorReset  = "\033[0m"
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorCyan   = "\033[36m"
)

// 目录信息结构
type DirInfo struct {
	Path      string
	FileCount int
	Files     []string
}

// 文件处理函数类型
// 参数：文件路径，返回：是否需要保存文件，错误信息
type FileHandler func(filePath string) (bool, error)

// 获取目录信息
func GetDirInfo(root string) (map[string]*DirInfo, error) {
	dirs := make(map[string]*DirInfo)

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("%s遍历目录时出错: %v%s\n", ColorRed, err, ColorReset)
			return err
		}

		if !d.IsDir() {
			relPath, err := filepath.Rel(root, path)
			if err != nil {
				fmt.Printf("%s获取相对路径失败 %s: %v%s\n", ColorRed, path, err, ColorReset)
				return err
			}

			dir := filepath.Dir(relPath)
			if _, exists := dirs[dir]; !exists {
				dirs[dir] = &DirInfo{
					Path:      dir,
					FileCount: 0,
					Files:     make([]string, 0),
				}
			}
			dirs[dir].FileCount++
			dirs[dir].Files = append(dirs[dir].Files, path)
		}
		return nil
	})

	return dirs, err
}

// 通用文件遍历器
// 参数：
//   - root: 要遍历的根目录
//   - handler: 自定义的文件处理函数
func ProcessFiles(root string, handler FileHandler) error {
	// 获取目录信息
	dirs, err := GetDirInfo(root)
	if err != nil {
		fmt.Printf("%s统计目录失败: %v%s\n", ColorRed, err, ColorReset)
		return err
	}

	// 处理每个目录
	for _, dirInfo := range dirs {
		fmt.Printf("\n%s开始处理目录%s: %s (%d 个文件)\n", ColorCyan, ColorReset, dirInfo.Path, dirInfo.FileCount)

		// 创建该目录的进度条
		bar := progressbar.NewOptions(dirInfo.FileCount,
			progressbar.OptionEnableColorCodes(true),
			progressbar.OptionShowBytes(false),
			progressbar.OptionSetWidth(15),
			progressbar.OptionSetTheme(progressbar.Theme{
				Saucer:        "[green]=[reset]",
				SaucerHead:    "[green]>[reset]",
				SaucerPadding: " ",
				BarStart:      "[",
				BarEnd:        "]",
			}),
		)

		// 处理目录中的每个文件
		for fileIdx, path := range dirInfo.Files {
			// 更新进度条描述
			bar.Describe(fmt.Sprintf("%s正在处理%s (%d/%d) %s", ColorCyan, ColorReset, fileIdx+1, dirInfo.FileCount, filepath.Base(path)))

			// 调用自定义处理函数
			needSave, err := handler(path)
			if err != nil {
				bar.Describe(fmt.Sprintf("%s处理失败%s %s", ColorRed, ColorReset, filepath.Base(path)))
				fmt.Printf("%s处理文件失败 %s: %v%s\n", ColorRed, path, err, ColorReset)
				return err
			}

			if needSave {
				bar.Describe(fmt.Sprintf("%s已更新%s %s", ColorBlue, ColorReset, filepath.Base(path)))
			} else {
				bar.Describe(fmt.Sprintf("%s无需更新%s %s", ColorYellow, ColorReset, filepath.Base(path)))
			}

			// 更新进度条
			_ = bar.Add(1)
		}

		_ = bar.Finish()
		fmt.Printf("%s目录处理完成%s: %s\n", ColorGreen, ColorReset, dirInfo.Path)
	}

	fmt.Printf("\n%s所有文件处理完成%s\n", ColorGreen, ColorReset)
	return nil
}

// 辅助函数：检查文件是否存在
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// 辅助函数：创建目录
func EnsureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// 辅助函数：读取文件内容
func ReadFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// 辅助函数：写入文件内容
func WriteFile(path string, data []byte) error {
	return os.WriteFile(path, data, 0644)
}
