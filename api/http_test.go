package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cookie = "ssid=eyJpdiI6IlwvWTl5K1Q5MkVKUFRoYkcrajRRakVRPT0iLCJ2YWx1ZSI6IkRINmhaQ09RUjhVRG1jc1BQQm5ESk52VXBLT2hsU0tGWW9lbXRVZlcyQ1d3cExyNjBMczhadHZQaWQyRzVyV3ZnVE5uSmdrd3lsTnJLUUcxN1F6QTVnPT0iLCJtYWMiOiI4NzdmYWFhMGU3NjQxYWEyZjFhNzM0NDNkNmFiODE3ZjQ2YWE0ZWQ1MGI1OWJlZGRlZjUyMzA0YzVmMTg3NDYyIn0%3D; 26efeb6279834e831e8659742d83d367=87dba676304ec89e18fbc47f443fc25d"

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
