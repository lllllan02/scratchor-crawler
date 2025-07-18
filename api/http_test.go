package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const cookie = "ssid=eyJpdiI6IjJraWtQZzAzb1VJTkloV3loVlUwZEE9PSIsInZhbHVlIjoieDl2S0dwK1N3MjRPaGM1aXhtY2RDdll2R2FuR0o5clZtcXl0THV4VlVvNkhrdVRFOHB0MzNvejZRV3BJekVhT01VaE1lQmhVNlUrWWlNTnpIRDVGa0E9PSIsIm1hYyI6IjRkZjE5NzExMDUwZjhjMDhjYTA1ZWRhZTRiNzdkYzNjOWMwYzc2M2Q1NzVhNDRhNzZmYzYzNmFmMDgxNGI1NDkifQ%3D%3D; 26efeb6279834e831e8659742d83d367=d6743583999a0cbf2984741b2a74a93a"

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
