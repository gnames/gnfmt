package gnfmt_test

import (
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/matryer/is"
)

func TestTimeString(t *testing.T) {
	is := is.New(t)
	tests := []struct {
		name string
		secs float64
		str  string
	}{
		{"zero", 0, "00:00:00"},
		{"sec", 30.3, "00:00:30"},
		{"min", 185.43432, "00:03:05"},
		{"hr", 7654.0003, "02:07:34"},
		{"d", 900_000, "10 days, 10:00:00"},
		{"dd", 90_000_160.1234, "1041 days, 16:02:40"},
		{"ddd", 90_000_000.1234, "1041 days, 16:00:00"},
	}

	for _, v := range tests {
		t.Run(v.name, func(_ *testing.T) {
			str := gnfmt.TimeString(v.secs)
			is.Equal(v.str, str)
		})
	}
}
