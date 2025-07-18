package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCat(t *testing.T) {
	links, err := GetCat(1, cookie)
	assert.NoError(t, err)
	fmt.Printf("links: %v\n", links)
}
