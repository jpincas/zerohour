package zerohour

import (
	"testing"
	"time"
)

const dateFormat = "Mon Jan 2 15:04:05 -0700 MST 2006"

// TODO: redo these 2 nasty tests, remove the parse and format

func TestPreviousUTC(t *testing.T) {
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
		pm := PreviousUTC(tn)
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

func TestRangeAdjust(t *testing.T) {
	tests := []struct {
		name             string
		from, to         time.Time
		expectedDaysDiff time.Duration
	}{
		{
			name:             "both blank",
			from:             time.Time{},
			to:               time.Time{},
			expectedDaysDiff: 1,
		},
		{
			name:             "same day, zero times",
			from:             time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			to:               time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			expectedDaysDiff: 1,
		},
		{
			name:             "same day, random times, from time is before to",
			from:             time.Date(2020, time.January, 2, 5, 12, 11, 75, time.UTC),
			to:               time.Date(2020, time.January, 2, 14, 24, 35, 99, time.UTC),
			expectedDaysDiff: 1,
		},
		{
			name:             "same day, random times, from time is after to",
			to:               time.Date(2020, time.January, 2, 5, 12, 11, 75, time.UTC),
			from:             time.Date(2020, time.January, 2, 14, 24, 35, 99, time.UTC),
			expectedDaysDiff: 1,
		},
		{
			name:             "expect 2 day interval, zero times",
			from:             time.Date(2020, time.January, 2, 0, 0, 0, 0, time.UTC),
			to:               time.Date(2020, time.January, 3, 0, 0, 0, 0, time.UTC),
			expectedDaysDiff: 2,
		},
		{
			name:             "expect 3 day interval, random times, from time is before to",
			from:             time.Date(2020, time.January, 2, 5, 12, 11, 75, time.UTC),
			to:               time.Date(2020, time.January, 4, 14, 24, 35, 99, time.UTC),
			expectedDaysDiff: 3,
		},
	}

	for _, test := range tests {
		diff := DaysDiffIgnoreTime(test.from, test.to)
		expectedDiff := test.expectedDaysDiff * 24 * time.Hour

		if diff != expectedDiff {
			t.Errorf("Testing %s.  Expected a diff of %v days, but got %v", test.name, expectedDiff, diff)
		}

	}

}

func TestStartOfDayInTimeZone(t *testing.T) {
	tests := []struct {
		name      string
		timeZone  string
		inputTime time.Time
		expected  time.Time
		wantErr   bool
	}{
		{
			name:      "timezone that is 4 hours behind UTC in April",
			timeZone:  "America/New_York",
			inputTime: time.Date(2023, 4, 27, 3, 30, 0, 0, time.UTC),
			expected:  time.Date(2023, 4, 27, 4, 0, 0, 0, time.UTC),
			wantErr:   false,
		},
		{
			name:      "timezone that is 9 hours ahead of UTC in April",
			timeZone:  "Asia/Tokyo",
			inputTime: time.Date(2023, 4, 27, 3, 30, 0, 0, time.UTC),
			expected:  time.Date(2023, 4, 26, 15, 0, 0, 0, time.UTC),
			wantErr:   false,
		},
		{
			name:      "invalid timezone",
			timeZone:  "Invalid/Timezone",
			inputTime: time.Date(2023, 4, 27, 3, 30, 0, 0, time.UTC),
			expected:  time.Now().UTC(),
			wantErr:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := StartOfDayInTimeZone(tt.inputTime, tt.timeZone)
			if (err != nil) != tt.wantErr {
				t.Errorf("StartOfDayInTimeZone() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && !got.Equal(tt.expected) {
				t.Errorf("StartOfDayInTimeZone() = %v, want %v", got, tt.expected)
			}
		})
	}
}
