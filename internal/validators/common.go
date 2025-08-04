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