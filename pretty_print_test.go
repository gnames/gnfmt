package gnfmt_test

import (
	"fmt"
	"math"
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/stretchr/testify/assert"
)

func TestPpr(t *testing.T) {
	assert := assert.New(t)
	type obj struct {
		A string
		B int
		C []int
		D []string
		E float32
	}

	o := obj{
		A: "one",
		B: 345,
		C: []int{1, 44},
		D: []string{"one", "two"},
		E: math.Pi,
	}
	res := gnfmt.Ppr(o)
	assert.Equal(1, 1)
	fmt.Println(res)
}
