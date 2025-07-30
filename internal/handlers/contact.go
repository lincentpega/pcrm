package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lincentpega/pcrm/internal/models"
	"github.com/lincentpega/pcrm/internal/repository"
	"github.com/lincentpega/pcrm/internal/templates"
)

type ContactHandler struct {
	contactRepo *repository.ContactRepository
	renderer    *templates.Renderer
}

func NewContactHandler(contactRepo *repository.ContactRepository, renderer *templates.Renderer) *ContactHandler {
	return &ContactHandler{
		contactRepo: contactRepo,
		renderer:    renderer,
	}
}

func (h *ContactHandler) NewContactForm(w http.ResponseWriter, r *http.Request) {
	personIDStr := r.PathValue("personId")
	personID, err := strconv.ParseInt(personIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	contactTypes, err := h.contactRepo.GetContactTypes()
	if err != nil {
		http.Error(w, "Failed to fetch contact types", http.StatusInternalServerError)
		return
	}

	data := struct {
		PersonID     int64
		ContactTypes []models.ContactType
		Form         *models.ContactForm
		IsEdit       bool
	}{
		PersonID:     personID,
		ContactTypes: contactTypes,
		Form:         &models.ContactForm{},
		IsEdit:       false,
	}

	h.renderer.RenderFragment(w, "contact_form", data)
}

func (h *ContactHandler) EditContactForm(w http.ResponseWriter, r *http.Request) {
	personIDStr := r.PathValue("personId")
	contactIDStr := r.PathValue("contactId")
	
	personID, err := strconv.ParseInt(personIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}
	
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	contact, err := h.contactRepo.GetByID(contactID)
	if err != nil {
		http.Error(w, "Contact not found", http.StatusNotFound)
		return
	}

	contactTypes, err := h.contactRepo.GetContactTypes()
	if err != nil {
		http.Error(w, "Failed to fetch contact types", http.StatusInternalServerError)
		return
	}

	data := struct {
		PersonID     int64
		ContactTypes []models.ContactType
		Form         *models.ContactForm
		Contact      *models.Contact
		IsEdit       bool
	}{
		PersonID:     personID,
		ContactTypes: contactTypes,
		Form:         contact.ToForm(),
		Contact:      contact,
		IsEdit:       true,
	}

	h.renderer.RenderFragment(w, "contact_form", data)
}

func (h *ContactHandler) CreateContact(w http.ResponseWriter, r *http.Request) {
	personIDStr := r.PathValue("personId")
	personID, err := strconv.ParseInt(personIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	contactTypeID, err := strconv.ParseInt(r.FormValue("contactTypeId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid contact type ID", http.StatusBadRequest)
		return
	}

	form := &models.ContactForm{
		ContactTypeID: contactTypeID,
		Content:       r.FormValue("content"),
	}

	contact := form.ToContact(personID)
	if err := h.contactRepo.Create(contact); err != nil {
		http.Error(w, "Failed to create contact", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/people/%d", personID))
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/people/%d", personID), http.StatusSeeOther)
}

func (h *ContactHandler) UpdateContact(w http.ResponseWriter, r *http.Request) {
	personIDStr := r.PathValue("personId")
	contactIDStr := r.PathValue("contactId")
	
	personID, err := strconv.ParseInt(personIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}
	
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	contactTypeID, err := strconv.ParseInt(r.FormValue("contactTypeId"), 10, 64)
	if err != nil {
		http.Error(w, "Invalid contact type ID", http.StatusBadRequest)
		return
	}

	contact := &models.Contact{
		ID:            contactID,
		PersonID:      personID,
		ContactTypeID: contactTypeID,
		Content:       r.FormValue("content"),
	}

	if err := h.contactRepo.Update(contact); err != nil {
		http.Error(w, "Failed to update contact", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/people/%d", personID))
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/people/%d", personID), http.StatusSeeOther)
}

func (h *ContactHandler) DeleteContact(w http.ResponseWriter, r *http.Request) {
	contactIDStr := r.PathValue("contactId")
	contactID, err := strconv.ParseInt(contactIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid contact ID", http.StatusBadRequest)
		return
	}

	if err := h.contactRepo.Delete(contactID); err != nil {
		http.Error(w, "Failed to delete contact", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.WriteHeader(http.StatusOK)
		return
	}

	personIDStr := r.PathValue("personId")
	if personIDStr != "" {
		http.Redirect(w, r, fmt.Sprintf("/people/%s", personIDStr), http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusOK)
	}
}