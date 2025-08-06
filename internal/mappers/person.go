package mappers

import (
	"github.com/lincentpega/pcrm/internal/dto"
	"github.com/lincentpega/pcrm/internal/models"
)

func PersonUpsertRequestToDomain(req *dto.PersonUpsertRequest) *models.Person {
	return &models.Person{
		FirstName:  req.FirstName,
		SecondName: req.SecondName,
		MiddleName: req.MiddleName,
	}
}

func PersonDomainToResponse(person *models.Person) dto.PersonInfoResponse {
	return dto.PersonInfoResponse{
		ID:         person.ID,
		FirstName:  person.FirstName,
		SecondName: person.SecondName,
		MiddleName: person.MiddleName,
		CreatedAt:  person.CreatedAt,
		UpdatedAt:  person.UpdatedAt,
	}
}

func ContactDomainToResponse(contact *models.Contact) dto.ContactResponse {
	return dto.ContactResponse{
		ID:        contact.ID,
		PersonID:  contact.PersonID,
		Content:   contact.Content,
		CreatedAt: contact.CreatedAt,
		UpdatedAt: contact.UpdatedAt,
		ContactType: dto.ContactTypeResponse{
			ID:        contact.ContactType.ID,
			Name:      contact.ContactType.Name,
			CreatedAt: contact.ContactType.CreatedAt,
		},
	}
}

func PersonWithContactsDomainToResponse(person *models.Person, contacts []models.Contact) dto.PersonWithContactsResponse {
	response := dto.PersonWithContactsResponse{
		PersonInfoResponse: PersonDomainToResponse(person),
		Contacts:           make([]dto.ContactResponse, len(contacts)),
	}

	for i, contact := range contacts {
		response.Contacts[i] = ContactDomainToResponse(&contact)
	}

	return response
}

