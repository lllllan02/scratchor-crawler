package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPaper(t *testing.T) {
	paper, err := client.GetPaper("https://tiku.scratchor.com/paper/view/zkq2vfoqpoiponm5")
	assert.NoError(t, err)
	fmt.Printf("paper: %+v\n", paper)
	for _, section := range paper.Sections {
		fmt.Printf("section: %+v\n", section)
	}
}
