package core

import (
	"log"
	"time"
)

func ParseDate(value *string) *time.Time {
	var dateValue *time.Time
	dateValue = nil
	if *value != "" {
		parsedDate, err := time.Parse("2006-01-02", *value)
		if err != nil {
			log.Fatal("invalid date value")
		}
		dateValue = &parsedDate
	}
	return dateValue
}
