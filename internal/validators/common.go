package validators

import (
	"errors"
	"strconv"

	"github.com/lincentpega/pcrm/internal/dto"
)

func ValidatePersonID(idStr string) (int64, error) {
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return 0, errors.New("invalid person ID")
	}
	return id, nil
}

func ValidatePersonUpsertRequest(req *dto.PersonUpsertRequest) error {
	if req.FirstName == "" {
		return errors.New("first name is required")
	}
	return nil
}

func ParsePaginationParams(pageStr, limitStr string) (int, int) {
	page := 1
	limit := 10

	if pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	if limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 100 {
			limit = l
		}
	}

	return page, limit
}

func ValidateBirthDateInfoRequest(req *dto.BirthDateInfoRequest) error {
	hasYear := req.BirthYear != nil
	hasMonth := req.BirthMonth != nil
	hasDay := req.BirthDay != nil
	hasApproximateAge := req.ApproximateAge != nil

	// Scenario 1: Exact date (year + month + day, no approximate age)
	if hasYear && hasMonth && hasDay && !hasApproximateAge {
		return validateExactDate(*req.BirthYear, *req.BirthMonth, *req.BirthDay)
	}

	// Scenario 2: Partial date (month + day only, no year or approximate age)
	if !hasYear && hasMonth && hasDay && !hasApproximateAge {
		return validatePartialDate(*req.BirthMonth, *req.BirthDay)
	}

	// Scenario 3: Approximate age only (no date fields)
	if !hasYear && !hasMonth && !hasDay && hasApproximateAge {
		return validateApproximateAge(*req.ApproximateAge)
	}

	// Scenario 4: Clear all (all fields nil)
	if !hasYear && !hasMonth && !hasDay && !hasApproximateAge {
		return nil // Valid - clearing birth date info
	}

	// Invalid combinations - provide specific error messages
	if hasApproximateAge && (hasYear || hasMonth || hasDay) {
		return errors.New("approximate age cannot be combined with date fields - use either approximate age alone or date fields without approximate age")
	}

	if hasYear && !hasMonth {
		return errors.New("birth year requires both month and day - provide complete date (year+month+day) or partial date (month+day only)")
	}

	if hasYear && hasMonth && !hasDay {
		return errors.New("birth year and month require day - provide complete date (year+month+day)")
	}

	if hasDay && !hasMonth {
		return errors.New("birth day requires month - provide either partial date (month+day) or complete date (year+month+day)")
	}

	// Fallback for any other invalid combinations
	return errors.New("invalid birth date combination - use exact date (year+month+day), partial date (month+day), approximate age only, or clear all fields")
}

func validateExactDate(year, month, day int) error {
	if year < 1900 || year > 2100 {
		return errors.New("birth year must be between 1900 and 2100")
	}
	if month < 1 || month > 12 {
		return errors.New("birth month must be between 1 and 12")
	}
	if day < 1 || day > 31 {
		return errors.New("birth day must be between 1 and 31")
	}

	// Check if the date is valid (handles leap years, month lengths)
	_, err := parseDate(year, month, day)
	if err != nil {
		return errors.New("invalid date")
	}

	return nil
}

func validatePartialDate(month, day int) error {
	if month < 1 || month > 12 {
		return errors.New("birth month must be between 1 and 12")
	}
	if day < 1 || day > 31 {
		return errors.New("birth day must be between 1 and 31")
	}

	// Check if month/day combination is valid (use leap year for February validation)
	_, err := parseDate(2024, month, day) // Use leap year to be permissive
	if err != nil {
		return errors.New("invalid month/day combination")
	}

	return nil
}

func validateApproximateAge(age int) error {
	if age < 0 || age > 150 {
		return errors.New("approximate age must be between 0 and 150")
	}
	return nil
}

func parseDate(year, month, day int) (bool, error) {
	// Check basic ranges
	if month < 1 || month > 12 {
		return false, errors.New("invalid month")
	}
	if day < 1 {
		return false, errors.New("invalid day")
	}

	// Days in each month
	daysInMonth := []int{31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31}
	
	// Check for leap year
	if isLeapYear(year) {
		daysInMonth[1] = 29
	}

	if day > daysInMonth[month-1] {
		return false, errors.New("day exceeds month limit")
	}

	return true, nil
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func ValidateConnectionSourceRequest(req *dto.ConnectionSourceRequest) error {
	if req.IntroducerPersonID != nil && req.IntroducerName != nil {
		return errors.New("cannot specify both introducer_person_id and introducer_name")
	}
	
	if req.WasIntroduced != nil && *req.WasIntroduced {
		if req.IntroducerPersonID == nil && req.IntroducerName == nil {
			return errors.New("when was_introduced is true, either introducer_person_id or introducer_name must be provided")
		}
	}
	
	if req.IntroducerPersonID != nil && *req.IntroducerPersonID <= 0 {
		return errors.New("introducer_person_id must be positive")
	}
	
	if req.IntroducerName != nil && *req.IntroducerName == "" {
		return errors.New("introducer_name cannot be empty")
	}
	
	return nil
}