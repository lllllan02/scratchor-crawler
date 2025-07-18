package crawler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cookie = "9c8adef9772544421fc21404af7ed346=522567bdcc2ae76b51518a4478f6cb80; ssid=eyJpdiI6ImpSVzI5b2dLNUxZaFJnZVhoNEc1ZFE9PSIsInZhbHVlIjoiK24rekhaRnoxMTduMHErS04yM1NUN1FpenVrUlZkMVJIVWFrU2QxWUdCSlVIQWFJR1FVc2VDQnhkVWQ5OGVqRjJGQjU5c2MwcytPbXZLYm5NS2NvY3c9PSIsIm1hYyI6ImExMmVlNThjZDNiOTk4NmMzYzFjYmRlMjk3Y2UwNjQ2OGM1NzIzYWRhYWU0NjY4YWIyN2Y3N2M5ODlhYTYwZGYifQ%3D%3D"

func TestGet(t *testing.T) {
	body, err := Get("https://tiku.scratchor.com/question/cat/1", cookie)
	assert.NoError(t, err)
	fmt.Printf("body: %v\n", body)
}

func TestPost(t *testing.T) {
	body, err := Post("https://tiku.scratchor.com/question/answer/abjyobpc2yhjb2ei", cookie, nil)
	assert.NoError(t, err)
	fmt.Printf("body: %v\n", body)
}
