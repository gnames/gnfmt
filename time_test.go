package gnfmt_test

import (
	"testing"

	"github.com/gnames/gnfmt"
	"github.com/stretchr/testify/assert"
)

func TestTimeString(t *testing.T) {
	tests := []struct {
		name string
		secs float64
		str  string
	}{
		{"zero", 0, "0sec"},
		{"sec", 30.3, "30sec"},
		{"min", 185.43432, "3m 5sec"},
		{"hr", 7654.0003, "2h 7m 34sec"},
		{"d", 88_000, "1d 26m 40sec"},
		{"dd", 90_000_160.1234, "1041d 16h 2m 40sec"},
		{"ddd", 90_000_000.1234, "1041d 16h 0m 0sec"},
	}

	for _, v := range tests {
		t.Run(v.name, func(_ *testing.T) {
			str := gnfmt.TimeString(v.secs)
			assert.Equal(t, v.str, str)
		})
	}
}
