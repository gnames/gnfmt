package gnfmt_test

import (
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/matryer/is"
)

func TestNewFormat(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		name   string
		input  string
		output gnfmt.Format
		errNil bool
	}{
		{"csv", "csv", gnfmt.CSV, true},
		{"compact", "compact", gnfmt.CompactJSON, true},
		{"pretty", "pretty", gnfmt.PrettyJSON, true},
		{"tsv", "tsv", gnfmt.TSV, true},
		{"bad", "bad", gnfmt.FormatNone, false},
	}
	for _, v := range tests {
		t.Run(v.name, func(_ *testing.T) {
			gf, err := gnfmt.NewFormat(v.input)
			is.Equal(v.output, gf)
			is.Equal(v.errNil, err == nil)
		})
	}
}

func TestString(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		name   string
		input  gnfmt.Format
		output string
	}{
		{"csv", gnfmt.CSV, "CSV"},
		{"tsv", gnfmt.TSV, "TSV"},
		{"compact", gnfmt.CompactJSON, "compact JSON"},
		{"pretty", gnfmt.PrettyJSON, "pretty JSON"},
	}
	for _, v := range tests {
		t.Run(v.name, func(_ *testing.T) {
			s := v.input.String()
			is.Equal(v.output, s)
		})
	}
}
