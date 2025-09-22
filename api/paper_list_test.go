package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPaperList(t *testing.T) {
	links, next, err := client.GetPaperList("https://tiku.scratchor.com/paper")
	assert.NoError(t, err)
	for _, link := range links {
		fmt.Printf("link: %s\n", link)
	}
	fmt.Printf("next: %s\n", next)
}
