package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCat(t *testing.T) {
	links, err := client.GetCat(1)
	assert.NoError(t, err)
	fmt.Printf("链接数量: %d\n", len(links))

	for i, link := range links {
		fmt.Printf("链接 %d: %s (题数: %d)\n", i+1, link.URL, link.Count)
	}
}
