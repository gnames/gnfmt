package gnfmt_test

import (
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
	assert.Equal("{\n  \"A\": \"one\",\n  \"B\": 345,\n  \"C\": [\n    1,\n    44\n  ],\n  \"D\": [\n    \"one\",\n    \"two\"\n  ],\n  \"E\": 3.1415927\n}", res)
}
