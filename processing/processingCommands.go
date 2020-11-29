package processing

import (
	"github.com/nicksnyder/go-i18n/i18n"
	"strconv"
	"strings"
	"time"
)

type timeLeft struct {
	years int
	months int
	days int
	hours int
}

const (
	hoursPerYear = 365 * 24
	hoursPerMonth = 30 * 24
	hoursPerDay = 24
)

func validateTimeLeft(parsedTime *timeLeft) (isSuccessfull bool, errorMessage string) {
	isSuccessfull = true
	return
}

func ParseBetTime(timeText string) (duration time.Duration, isSuccessfull bool, errorMessage string) {	
	// possible date formats DD.MM.YY DD.MM.YYYY YYYYMMDD YYMMDD
	// possible duration formats 1y2m3d4h 2m 2y6m ...

	if timeText == "" {
		errorMessage = "Empty time format is not allowed"
		isSuccessfull = false
		return
	}

	if strings.ContainsAny(timeText, "ymdh") {
		var parsedTime timeLeft
		var daysOrHoursParsed bool
		
		textToParse := timeText
		for {
			nextIdx := strings.IndexAny(textToParse, "ymdh")
			if nextIdx == -1 {
				break
			}
			value, err := strconv.Atoi(textToParse[:nextIdx])
			if err != nil {
				errorMessage = "Invalid numeric value: " + textToParse[:nextIdx]
				isSuccessfull = false
				return
			}

			switch textToParse[nextIdx] {
			case 'y':
				parsedTime.years = value
			case 'm':
				if daysOrHoursParsed {
				errorMessage = "Invalid order for months value (there are no minutes)"
					isSuccessfull = false
					return
				}
				parsedTime.months = value
			case 'd':
				parsedTime.days = value
				daysOrHoursParsed = true
			case 'h':
				parsedTime.hours = value
				daysOrHoursParsed = true
			}

			textToParse = textToParse[nextIdx+1:]
		}

		isSuccessfull, errorMessage = validateTimeLeft(&parsedTime)

		duration = time.Duration(
			time.Duration(parsedTime.years * hoursPerYear) * time.Hour +
			time.Duration(parsedTime.months * hoursPerMonth) * time.Hour +
			time.Duration(parsedTime.days * hoursPerDay) * time.Hour +
			time.Duration(parsedTime.hours) * time.Hour)
	} else {
		
	}
	return
}

func getTimeLeftFromDuration(duration time.Duration) (resultTime timeLeft) {
	hours := int(duration.Hours())

	resultTime.years = hours / hoursPerYear
	hours = hours - resultTime.years * hoursPerYear
	resultTime.months = hours / hoursPerMonth
	hours = hours - resultTime.months * hoursPerMonth
	resultTime.days = hours / hoursPerDay
	resultTime.hours = hours - resultTime.days * hoursPerDay
	return
}

func getTimeFormat(parsedTime *timeLeft) string {
	if parsedTime.years > 0 {
		if parsedTime.months > 0 {
			return "time_left_years"
		} else {
			return "time_left_years_nomonths"
		}
	} else if parsedTime.months > 0 {
		if parsedTime.days > 0 {
			return "time_left_months"
		} else {
			return "time_left_months_nodays"
		}
	} else if parsedTime.days > 0 {
		if parsedTime.hours > 0 {
			return "time_left_days"
		} else {
			return "time_left_days_nohours"
		}
	} else if parsedTime.hours > 0 {
		return "time_left_hours"
	} else {
		return "time_left_lessthanhour"
	}
}

func GetBetDurationText(duration time.Duration, trans i18n.TranslateFunc) string {
	parsedTime := getTimeLeftFromDuration(duration)

	rulesData := map[string]interface{}{
		"Years": trans("years", parsedTime.years),
		"Months": trans("months", parsedTime.months),
		"Days": trans("days", parsedTime.days),
		"Hours": trans("hours", parsedTime.hours),
	}

	return trans(getTimeFormat(&parsedTime), rulesData)
}
