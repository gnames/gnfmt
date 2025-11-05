package gnfmt

import (
	"fmt"
	"strings"
)

func TimeString(secs float64) string {
	type duration struct {
		days    int
		hours   int
		minutes int
		seconds int
	}

	var (
		day     = 86400
		hour    = 3600
		min     = 60
		secsInt = int(secs)
		dur     = duration{}
		resid   = secsInt
	)

	dur.days = resid / day
	if dur.days > 0 {
		resid = resid % day
	}

	dur.hours = resid / hour
	if dur.hours > 0 {
		resid = resid % hour
	}

	dur.minutes = resid / min
	if dur.minutes > 0 {
		resid = resid % min
	}

	dur.seconds = resid
	res := ""

	if dur.days > 0 {
		res += fmt.Sprintf("%dd ", dur.days)
	}

	if dur.hours > 0 {
		res += fmt.Sprintf("%dh ", dur.hours)
	}

	if dur.hours > 0 || dur.minutes > 0 {
		res += fmt.Sprintf("%dm ", dur.minutes)
	}

	res += fmt.Sprintf("%dsec ", dur.seconds)

	return strings.TrimSpace(res)
}
