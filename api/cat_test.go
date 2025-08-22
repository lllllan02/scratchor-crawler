package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCat(t *testing.T) {
	cat, err := client.GetCat(1)
	assert.NoError(t, err)

	fmt.Printf("分类: %s\n", cat.Title)
	for _, link := range cat.Links {
		fmt.Printf("%s: %s (题数: %d)\n", link.Title, link.URL, link.Count)
	}
}
