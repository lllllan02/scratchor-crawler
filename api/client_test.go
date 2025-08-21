package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client, _ = NewClient("0368c6a3166608f49a794ade04580c8c=314d2208c434ca1420c72f42b8de2695; Hm_lvt_22880f1d42e3788b94e7aa1361ea923a=1755757549; HMACCOUNT=4D5E77D5BB397FC9; ssid=eyJpdiI6IlU3MUJ4YmVZb0djUTJmUnVVOFpTQ2c9PSIsInZhbHVlIjoiYnFiZk5nZHBxb0ZyN3ozNWFXKzl2b2tFXC83NHppZUwyVFlhUmJsbFczMHR2RFJSSEpCTTZnOUVUa1Q3NFVIWWRPRGk0bFFKV2hVNXkzNFhockFzK2J3PT0iLCJtYWMiOiJkNWUxMjBjYmU4M2M2NGI3ZTRjNjNhZWZiOWZlZmVjM2NlNWY0YzcyNDYyNWE0ODJjMzYyNGNhYmRhY2ZiYzNhIn0%3D; Hm_lpvt_22880f1d42e3788b94e7aa1361ea923a=1755757751")

func TestImage(t *testing.T) {
	body, err := client.Image("https://imgtk.scratchor.com/data/image/2022/11/02/86228_619t_2689.png")
	assert.NoError(t, err)
	os.WriteFile("test.png", body, 0644)
}
