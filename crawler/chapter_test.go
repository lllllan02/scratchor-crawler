package crawler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetChapter(t *testing.T) {
	links, next, err := GetChapter("https://tiku.scratchor.com/question/cat/1/list?chapterId=1&page=1", cookie)
	assert.NoError(t, err)
	fmt.Printf("links: %v\n", links)
	fmt.Printf("next: %v\n", next)
}
