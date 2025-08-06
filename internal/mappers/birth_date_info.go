package mappers

import (
	"time"

	"github.com/lincentpega/pcrm/internal/dto"
	"github.com/lincentpega/pcrm/internal/models"
)

func BirthDateInfoRequestToDomain(personID int64, req *dto.BirthDateInfoRequest) *models.BirthDateInfo {
	var approximateAgeUpdatedAt *time.Time
	if req.ApproximateAge != nil {
		now := time.Now()
		approximateAgeUpdatedAt = &now
	}

	return &models.BirthDateInfo{
		PersonID:                personID,
		BirthYear:               req.BirthYear,
		BirthMonth:              req.BirthMonth,
		BirthDay:                req.BirthDay,
		ApproximateAge:          req.ApproximateAge,
		ApproximateAgeUpdatedAt: approximateAgeUpdatedAt,
	}
}

func BirthDateInfoDomainToResponse(birthDateInfo *models.BirthDateInfo) dto.BirthDateInfoResponse {
	return dto.BirthDateInfoResponse{
		ID:                      birthDateInfo.ID,
		PersonID:                birthDateInfo.PersonID,
		BirthYear:               birthDateInfo.BirthYear,
		BirthMonth:              birthDateInfo.BirthMonth,
		BirthDay:                birthDateInfo.BirthDay,
		ApproximateAge:          birthDateInfo.ApproximateAge,
		ApproximateAgeUpdatedAt: birthDateInfo.ApproximateAgeUpdatedAt,
		CreatedAt:               birthDateInfo.CreatedAt,
		UpdatedAt:               birthDateInfo.UpdatedAt,
	}
}