package zerohour

import "time"

// EndOfDay returns a time equivalent to the very last microsend of the day corresponding
// to the supplied time
func EndOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 23, 59, 59, 9999, t.Location())
}

// StartOfDay returns a time equivalent to the very first moment of the day corresponding
// to the supplied time - i.e. the 'zero hour' (00:00)
func StartOfDay(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
}

// StartOfTodayInTimeZone returns the UTC time corresponding to the start of the day
// that is currently active in a given timezone
func StartOfTodayInTimeZone(tz string) time.Time {
	loc, err := time.LoadLocation(tz)
	if err != nil {
		return time.Now().UTC(), err
	}

	t := time.Now().In(loc)
	return StartOfDay(t).UTC()
}

// FromToIgnoreTime converts a date range ('from' and 'to') dates into their
// time-agnostic times (i.e. it discards time data). Useful to adjusting date ranges
// where you want to include everything from the very start of the 'from' date right
// through to the end of the 'to' date.
func FromToIgnoreTime(from, to time.Time) (time.Time, time.Time) {
	return StartOfDay(from), EndOfDay(to)
}

// DaysDiffIngoreTime works out the number of 24 hour periods between the 'to' and 'from' times adjusted
// using FromToIgnoreTime, i.e. given times corresponding to the same day, the result will be 1 irrespective
// of the hours/minutes/seconds
func DaysDiffIgnoreTime(from, to time.Time) time.Duration {
	// We 'expect' the range adjuster to basically give a 'day' difference
	// given the same moment in time.  To test for this expectation, what I'm doing is rounding the
	// diff to the closest hour (because otherwise its a microsend under 24 hours)
	adjustedFrom, adjustedTo := FromToIgnoreTime(from, to)
	return adjustedTo.Sub(adjustedFrom).Round(time.Hour)
}

// Previous returns the time of the most recently passed zero hour (00:00) for any given time.
// In the infinitely rare case of the given time being precisely midnight (zero hour), it will
// return the zero hour 24 hours earlier.
func Previous(t time.Time) time.Time {
	if t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 {
		t = t.Add(time.Second * -1)
	}

	return StartOfDay(t)
}

// Previous returns the time of the most recently passed UTC zero hour (00:00) for any given time.
// NOTE that this doesn't just change the timezone of the result to UTC, but actually works out
// when the last UTC midnight was compared to the given time.  The result IS given back in UTC, so can be
// changed back to the original timezone if required.
// In the infinitely rare case of the given time being precisely midnight (zero hour), it will
// return the zero hour 24 hours earlier.
func PreviousUTC(t time.Time) time.Time {
	t = t.UTC()
	return Previous(t)
}

// PreviousSpecificDay returns the time of the most recently passed zero hour (00:00)
// but for a specific day, e.g. Sunday
func PreviousSpecificDay(t time.Time, targetDay time.Weekday) time.Time {
	previousZeroHour := Previous(t)
	dayOfPreviousZeroHour := previousZeroHour.Weekday()

	// If it happens to be the target day, just return it
	if dayOfPreviousZeroHour == targetDay {
		return previousZeroHour
	}

	// otherwise deduct the number of days away from the target day that we currently are
	return previousZeroHour.Add(-(time.Hour * 24 * time.Duration(dayOfPreviousZeroHour-targetDay)))
}

// PreviousSpecificDay returns the time of the most recently passed zero hour (00:00)
// but for a specific day, e.g. Sunday
// NOTE that this doesn't just change the timezone of the result to UTC, but actually works out
// when the last UTC midnight was compared to the given time.  The result IS given back in UTC, so can be
// changed back to the original timezone if required.
func PreviousSpecificDayUTC(t time.Time, targetDay time.Weekday) time.Time {
	t = t.UTC()
	return PreviousSpecificDay(t, targetDay)
}
