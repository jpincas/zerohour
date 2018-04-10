package zerohour

import "time"

// Previous returns the time of the most recently passed zero hour (00:00)
func Previous(now time.Time) time.Time {
	now = now.UTC()
	d := now.Day()
	// infintately rare case of it being exactly midnight
	// decrement the day by 1
	if now.Hour() == 0 && now.Minute() == 0 && now.Second() == 0 {
		d--
	}

	return time.Date(now.Year(), now.Month(), d, 0, 0, 0, 0, now.Location())
}

// PreviousSpecificDay returns the time of the most recently passed zero hour (00:00)
// but for a specific day, e.g. Sunday
func PreviousSpecificDay(now time.Time, day time.Weekday) time.Time {
	now = now.UTC()

	date := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	d := date.Weekday()

	// If it happens to be the target day, just return it
	if d == day {
		return date
	}

	// otherwise deduct the number of days away from the target day that we currently are
	return date.Add(-(time.Hour * 24 * time.Duration(d-day)))
}
