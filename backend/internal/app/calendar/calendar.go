package calendar

import "time"

func GetWeekdayCurrentWeek(weekday time.Weekday) time.Time {
	return GetWeekdayRelative(time.Now(), weekday)
}

func GetWeekdayRelative(now time.Time, weekday time.Weekday) time.Time {
	currentWeekday := int(now.Weekday())
	targetWeekday := int(weekday)

	if targetWeekday == int(time.Sunday) {
		targetWeekday = 7
	}

	if currentWeekday == int(time.Sunday) {
		currentWeekday = 7
	}

	daysUntilWeekday := targetWeekday - currentWeekday

	return now.AddDate(0, 0, daysUntilWeekday)
}
