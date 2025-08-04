package services

import (
	"fmt"
	"time"

	"github.com/lincentpega/pcrm/internal/models"
)

func FormatBirthDateInfo(person *models.Person) string {
	if person.BirthYear != nil && person.BirthMonth != nil && person.BirthDay != nil {
		return fmt.Sprintf("%d %s %d", 
			*person.BirthDay,
			getMonthName(*person.BirthMonth), 
			*person.BirthYear)
	}

	if person.BirthMonth != nil && person.BirthDay != nil {
		return fmt.Sprintf("%d %s", 
			*person.BirthDay,
			getMonthName(*person.BirthMonth))
	}

	if person.ApproximateAge != nil && person.ApproximateAgeUpdatedAt != nil {
		yearsPassed := calculateYearsPassed(*person.ApproximateAgeUpdatedAt, time.Now())
		currentApproximateAge := *person.ApproximateAge + yearsPassed
		return fmt.Sprintf("~%d years old", currentApproximateAge)
	}

	return "Age unknown"
}

func calculateYearsPassed(from, to time.Time) int {
	years := to.Year() - from.Year()
	
	if to.Month() < from.Month() || (to.Month() == from.Month() && to.Day() < from.Day()) {
		years--
	}
	
	return years
}

func getMonthName(month int) string {
	months := []string{
		"", "January", "February", "March", "April", "May", "June",
		"July", "August", "September", "October", "November", "December",
	}
	
	return months[month]
}