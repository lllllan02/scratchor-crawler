package utils

import "strings"

func CleanFileName(name string) string {
	return strings.ReplaceAll(name, "/", "_")
}
