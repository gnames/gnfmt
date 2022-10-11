package gnfmt

import "fmt"

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
	switch dur.days {
	case 0:
	case 1:
		res += fmt.Sprintf("%d day ", dur.days)
	default:
		res += fmt.Sprintf("%d days ", dur.days)
	}
	res += fmt.Sprintf("%02d:%02d:%02d", dur.hours, dur.minutes, dur.seconds)
	return res
}
