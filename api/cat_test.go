package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCat(t *testing.T) {
	links, err := client.GetCat(1)
	assert.NoError(t, err)
	fmt.Printf("links: %v\n", links)
}
