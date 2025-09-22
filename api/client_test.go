package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client, _ = NewClient("HMACCOUNT=4D5E77D5BB397FC9; Hm_lvt_22880f1d42e3788b94e7aa1361ea923a=1758358879; a0e0cf20f929d48c545b6794f503b8d2=14ea9ecc4ecea3c94148df41959a843d; ssid=eyJpdiI6ImdxUGpWRkhGMWdyZmwzXC9XSElhKzlRPT0iLCJ2YWx1ZSI6ImNqemx6Y2poNFNMU2tLWE1wc0VhSFVMSFVXOHM2S090bUtLK3JJbTI4bjFyYnpJMk5BXC9vUzZYXC9MVU9pRDNaNGIralJYMDAyU1h3OVpWa1VJOEdBY1E9PSIsIm1hYyI6ImUzZDk0NDU1NzUyZmJjYjU1N2EyYzEzMzA0NmM5ZGJkYzFkMGYwYTA2NTFhNWEzMzI4MmY2OGI2NmFkZTI1MmEifQ%3D%3D; Hm_lpvt_22880f1d42e3788b94e7aa1361ea923a=1758502503")

func TestImage(t *testing.T) {
	body, err := client.Image("https://imgtk.scratchor.com/data/image/2022/11/02/86228_619t_2689.png")
	assert.NoError(t, err)
	os.WriteFile("test.png", body, 0644)
}
