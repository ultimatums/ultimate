package durutil

import (
	"fmt"
	"regexp"
	"strconv"
	"time"
)

var durationPT = regexp.MustCompile("^([0-9]+)([ywdhms]+)$")

func DurationToString(duration time.Duration) string {
	seconds := int64(duration / time.Second)
	factors := map[string]int64{
		"y": 60 * 60 * 24 * 365,
		"w": 60 * 60 * 24 * 7,
		"d": 60 * 60 * 24,
		"h": 60 * 60,
		"m": 60,
		"s": 1,
	}
	unit := "s"
	switch int64(0) {
	case seconds % factors["y"]:
		unit = "y"
	case seconds % factors["w"]:
		unit = "w"
	case seconds % factors["d"]:
		unit = "d"
	case seconds % factors["h"]:
		unit = "h"
	case seconds % factors["m"]:
		unit = "m"
	}
	return fmt.Sprintf("%v%v", seconds/factors[unit], unit)
}

func StringToDuration(durationStr string) (duration time.Duration, err error) {
	matches := durationPT.FindStringSubmatch(durationStr)
	if len(matches) != 3 {
		err = fmt.Errorf("not a valid duration string: %q", durationStr)
		return
	}
	durationSeconds, _ := strconv.Atoi(matches[1])
	duration = time.Duration(durationSeconds) * time.Second
	unit := matches[2]
	switch unit {
	case "y":
		duration *= 60 * 60 * 24 * 365
	case "w":
		duration *= 60 * 60 * 24 * 7
	case "d":
		duration *= 60 * 60 * 24
	case "h":
		duration *= 60 * 60
	case "m":
		duration *= 60
	case "s":
		duration *= 1
	default:
		return 0, fmt.Errorf("invalid time unit in duration string: %q", unit)
	}
	return
}
