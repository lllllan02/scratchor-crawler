package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAnswer(t *testing.T) {
	answer, err := client.GetAnswer("abjyobpc2yhjb2ei")
	assert.NoError(t, err)
	assert.Equal(t, "D", answer)

	_, err = client.GetAnswer("hrslpuuqoqvaytgk")
	assert.Equal(t, ErrDailyLimitExceeded, err)
}
