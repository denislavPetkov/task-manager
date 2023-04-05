package controller

import "time"

const (
	inputDateFormat   = "2006-01-02"
	desiredDateFormat = "02-Jan-2006"
)

func getFormattedDeadline(deadline string) (string, error) {
	deadlineTime, err := time.Parse(inputDateFormat, deadline)
	if err != nil {
		return "", err
	}

	formattedDeadline := deadlineTime.Format(desiredDateFormat)

	return formattedDeadline, nil
}
