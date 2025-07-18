package api

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAnswer(t *testing.T) {
	answer, err := GetAnswer("abjyobpc2yhjb2ei", cookie)
	assert.NoError(t, err)
	assert.Equal(t, "D", answer)

	_, err = GetAnswer("hrslpuuqoqvaytgk", cookie)
	assert.Equal(t, ErrDailyLimitExceeded, err)
}
