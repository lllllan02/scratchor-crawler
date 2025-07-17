package crawler

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cookie = "af3eb5957ff8f6344db30bdde1473af2=75ff0ededc8526e42ec5f00103ed3bec; ssid=eyJpdiI6ImtvZ1E5QnNcL1htYmEyc3hzZkd3YjdBPT0iLCJ2YWx1ZSI6IlArZ0hQWWd0aWt4c1o1aHNoUDg4NzRWZ3NwZElYQVNLcmlrNTREUUN4aVRQRmQ4ZFdzdjVOMFVHUU9xVTJiK2ZcL05qaUNrY3cxZkVzcjFFT05TQldmdz09IiwibWFjIjoiMDY5MGU2ODZlYjcwZDU3YmIwNTUzMWIyNGYwNDY5OWQzNjE5ZmE0MTNlNDJlNGU4YmU0NGU5YTk2ZTM3NGNmMCJ9; 3b8e330e9dc4ceca2deb76a31d27f816=2a9a15f1a7fa9fb2c764febe850b2693"

func TestGet(t *testing.T) {
	body, err := Get("https://tiku.scratchor.com/question/cat/1", cookie)
	assert.NoError(t, err)
	fmt.Printf("body: %v\n", body)
}
