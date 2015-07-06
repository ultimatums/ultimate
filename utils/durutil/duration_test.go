package durutil

import (
	"testing"
	"time"
)

func TestStringToDuration(t *testing.T) {

	testCases := []string{
		"7h",
		"2y",
		"3w",
		"4d",
		"5h",
		"7m",
	}

	for _, testCase := range testCases {
		dur, err := StringToDuration(testCase)
		if err != nil {
			t.Error(err)
		}
		if testCase != DurationToString(dur) {
			t.Error("Duration to string failed.")
		}
	}

}

func Test(t *testing.T) {
	dur, err := time.ParseDuration("3h")
	if err != nil {
		t.Error(err)
	}
	t.Log(dur)
}
