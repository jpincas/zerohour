package zerohour

import (
	"testing"
	"time"
)

const dateFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"

func TestPreviousZeroHour(t *testing.T) {
	testCases := []struct {
		timeNow, expectedPreviousMidnight string
	}{
		{"Mon Jan 2 15:04:05 -0700 MST 2006", "Mon Jan 2 00:00:00 +0000 UTC 2006"},
		// take a time zone 7 hours ahead of UTC where it is currently 6am
		// therefore UTC time is 23:00 the previous day
		// therefore previous midnight should be 00:00 on that previous day
		{"Mon Jan 2 06:00:00 +0700 ICT 2006", "Mon Jan 1 00:00:00 +0000 UTC 2006"},
		{"Mon Jan 2 00:00:05 +0000 UTC 2006", "Mon Jan 2 00:00:00 +0000 UTC 2006"},
		{"Mon Mar 15 00:00:01 +0000 UTC 2006", "Mon Mar 15 00:00:00 +0000 UTC 2006"},
		{"Mon Mar 15 23:59:59 +0000 UTC 2018", "Mon Mar 15 00:00:00 +0000 UTC 2018"},
		{"Mon Mar 16 00:00:00 +0000 UTC 2018", "Mon Mar 15 00:00:00 +0000 UTC 2018"},
	}

	for _, testCase := range testCases {
		// parse the times
		tn, _ := time.Parse(dateFormat, testCase.timeNow)
		epm, _ := time.Parse(dateFormat, testCase.expectedPreviousMidnight)

		// run 'previous midnight' and stringify it
		pm := Previous(tn)
		pmString := pm.Format(dateFormat)

		if !pm.Equal(epm) {
			t.Errorf("Previous 0-hour calculation incorrect. Expected %s; Got %s", testCase.expectedPreviousMidnight, pmString)
		}
	}

}

func TestPreviousSaturdayMidnight(t *testing.T) {
	testCases := []struct {
		timeNow, expectedPreviousMidnight string
	}{
		{"Sun Apr 8 09:00:00 +0000 UTC 2018", "Sun Apr 8 00:00:00 +0000 UTC 2018"},
		// even if its only 1 second past midnight on saturday night
		{"Sun Apr 8 00:00:01 +0000 UTC 2018", "Sun Apr 8 00:00:00 +0000 UTC 2018"},
		{"Mon Apr 9 09:00:00 +0000 UTC 2018", "Sun Apr 8 00:00:00 +0000 UTC 2018"},
		{"Wed Apr 11 23:59:59 +0000 UTC 2018", "Sun Apr 8 00:00:00 +0000 UTC 2018"},
		{"Wed Apr 11 23:59:59 +0000 UTC 2018", "Sun Apr 8 00:00:00 +0000 UTC 2018"},
	}

	for _, testCase := range testCases {
		// parse the times
		tn, _ := time.Parse(dateFormat, testCase.timeNow)
		epm, _ := time.Parse(dateFormat, testCase.expectedPreviousMidnight)

		// run 'previous midnight' and stringify it
		pm := PreviousSpecificDay(tn, time.Sunday)
		pmString := pm.Format(dateFormat)

		if !pm.Equal(epm) {
			t.Errorf("Previous saturday midnight calculation incorrect. Expected %s; Got %s", testCase.expectedPreviousMidnight, pmString)
		}
	}

}
