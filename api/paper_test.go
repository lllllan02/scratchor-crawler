package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPaper(t *testing.T) {
	links, next, err := client.GetPaper("https://tiku.scratchor.com/paper")
	assert.NoError(t, err)
	for _, link := range links {
		fmt.Printf("link: %s\n", link)
	}
	fmt.Printf("next: %s\n", next)
}
