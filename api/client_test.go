package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client, _ = NewClient("26efeb6279834e831e8659742d83d367=d6743583999a0cbf2984741b2a74a93a; ssid=eyJpdiI6Im9PbWRMVUxLT3N2Zjk3OVZLMjJHbEE9PSIsInZhbHVlIjoiSmdSM1lqYlN1ZGE1VUJ5Zkl5NEpYY290d2xYWUNEamkyZWR4TnlHT3VHQ2ZiR2dTTTd3WWEwQW9HK3k2R3pXU3hYWkxHZjFUc0RyNTBoTUttYjExVGc9PSIsIm1hYyI6ImNiM2Q1MDg0ZWQwNmJjZjcwMmY3NWU2NDgxOTlhMmVkMDAxZjJiNzE3YmRhMGNmZmMwOTIyNGZmYjc5NGY0ZDYifQ%3D%3D")

func TestImage(t *testing.T) {
	body, err := client.Image("https://imgtk.scratchor.com/data/image/2022/11/02/86228_619t_2689.png")
	assert.NoError(t, err)
	os.WriteFile("test.png", body, 0644)
}
