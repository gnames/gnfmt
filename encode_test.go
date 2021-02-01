package gnfmt_test

import (
	"fmt"
	"testing"

	"github.com/matryer/is"
	"github.com/stretchr/testify/assert"

	. "github.com/gnames/gnfmt"
)

type version struct {
	Version string
	Build   string
}

func TestEncodeDecode(t *testing.T) {
	encs := []Encoder{
		GNgob{},
		GNjson{},
	}
	for _, e := range encs {
		obj := version{
			Version: "v10.10.10",
			Build:   "today",
		}
		res, err := e.Encode(obj)
		assert.Nil(t, err)
		var ver version
		err = e.Decode(res, &ver)
		assert.Nil(t, err)
		assert.Equal(t, ver.Version, "v10.10.10")
		assert.Equal(t, ver.Build, "today")
	}
}

func TestOutput(t *testing.T) {
	is := is.New(t)

	tests := []struct {
		name   string
		format Format
		input  interface{}
		output string
	}{
		{
			name:   "pretty",
			format: PrettyJSON,
			input: version{
				Version: "v10.10.10",
				Build:   "today",
			},
			output: `{
  "Version": "v10.10.10",
  "Build": "today"
}`,
		},
		{
			name:   "compact",
			format: CompactJSON,
			input: version{
				Version: "v10.10.10",
				Build:   "today",
			},
			output: `{"Version":"v10.10.10","Build":"today"}`,
		},
		{
			name:   "csv",
			format: CSV,
			input: version{
				Version: "v10.10.10",
				Build:   "today",
			},
			output: "",
		},
	}

	for _, v := range tests {
		t.Run(v.name, func(t *testing.T) {
			enc := GNjson{}
			s := enc.Output(v.input, v.format)
			is.Equal(v.output, s)
		})
	}
}

func Example() {
	var enc Encoder
	var err error
	enc = GNjson{Pretty: true}
	ary1 := []int{1, 2, 3}
	jsonRes, err := enc.Encode(ary1)
	if err != nil {
		panic(err)
	}
	var ary2 []int
	err = enc.Decode(jsonRes, &ary2)
	if err != nil {
		panic(err)
	}
	fmt.Println(ary1[0] == ary2[0])
	// Output: true
}
