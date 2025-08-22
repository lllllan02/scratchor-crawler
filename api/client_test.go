package api

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var client, _ = NewClient("Hm_lvt_22880f1d42e3788b94e7aa1361ea923a=1755757549; HMACCOUNT=4D5E77D5BB397FC9; 333decb516a63de949aa73f356ff0515=51967059785189c5bd710720adcf8afb; ssid=eyJpdiI6ImFLeDVPd3AyVVExVlR2VzhpRGVTdXc9PSIsInZhbHVlIjoiNXF3Wk10TE40OUl2OEhCY252cFVjT0dZVmd2eUdrRlNiXC85QjY4Q3Ztdk5HWTFLaTlVTjJSS1hZcUxrd1RDNXZSZ3BqOVN2XC8yc1NSTFBtdlBBTjdpdz09IiwibWFjIjoiMDU3MzllNjBmMmZlYzFjZTQ5N2EwMmI1NjdkNmE3N2RkNTY5ZDM2YjAyOGY2ZjhlZDg2YWM3MzkzODdiMjMwMiJ9; Hm_lpvt_22880f1d42e3788b94e7aa1361ea923a=1755854532")

func TestImage(t *testing.T) {
	body, err := client.Image("https://imgtk.scratchor.com/data/image/2022/11/02/86228_619t_2689.png")
	assert.NoError(t, err)
	os.WriteFile("test.png", body, 0644)
}
