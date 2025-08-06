package mappers

import (
	"github.com/lincentpega/pcrm/internal/dto"
	"github.com/lincentpega/pcrm/internal/models"
)

func ConnectionSourceRequestToDomain(personID int64, req *dto.ConnectionSourceRequest) *models.ConnectionSource {
	return &models.ConnectionSource{
		PersonID:           personID,
		MeetingStory:       req.MeetingStory,
		MeetingTimestamp:   req.MeetingTimestamp,
		WasIntroduced:      req.WasIntroduced,
		IntroducerPersonID: req.IntroducerPersonID,
		IntroducerName:     req.IntroducerName,
	}
}

func ConnectionSourceDomainToResponse(connectionSource *models.ConnectionSource) dto.ConnectionSourceResponse {
	return dto.ConnectionSourceResponse{
		ID:                 connectionSource.ID,
		PersonID:           connectionSource.PersonID,
		MeetingStory:       connectionSource.MeetingStory,
		MeetingTimestamp:   connectionSource.MeetingTimestamp,
		WasIntroduced:      connectionSource.WasIntroduced,
		IntroducerPersonID: connectionSource.IntroducerPersonID,
		IntroducerName:     connectionSource.IntroducerName,
		CreatedAt:          connectionSource.CreatedAt,
		UpdatedAt:          connectionSource.UpdatedAt,
	}
}