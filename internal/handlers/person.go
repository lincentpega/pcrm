package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/lincentpega/pcrm/internal/models"
	"github.com/lincentpega/pcrm/internal/repository"
	"github.com/lincentpega/pcrm/internal/templates"
)

type PersonHandler struct {
	repo        *repository.PersonRepository
	contactRepo *repository.ContactRepository
	renderer    *templates.Renderer
}

func NewPersonHandler(repo *repository.PersonRepository, contactRepo *repository.ContactRepository, renderer *templates.Renderer) *PersonHandler {
	return &PersonHandler{
		repo:        repo,
		contactRepo: contactRepo,
		renderer:    renderer,
	}
}

func (h *PersonHandler) ListPeople(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 10
	
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if p, err := strconv.Atoi(pageStr); err == nil && p > 0 {
			page = p
		}
	}

	people, err := h.repo.GetPaginated(page, limit)
	if err != nil {
		http.Error(w, "Failed to fetch people", http.StatusInternalServerError)
		return
	}

	totalCount, err := h.repo.GetTotalCount()
	if err != nil {
		http.Error(w, "Failed to get total count", http.StatusInternalServerError)
		return
	}

	totalPages := (totalCount + limit - 1) / limit

	data := struct {
		Title       string
		People      []models.Person
		CurrentPage int
		TotalPages  int
		HasPrev     bool
		HasNext     bool
		PrevPage    int
		NextPage    int
	}{
		Title:       "People",
		People:      people,
		CurrentPage: page,
		TotalPages:  totalPages,
		HasPrev:     page > 1,
		HasNext:     page < totalPages,
		PrevPage:    page - 1,
		NextPage:    page + 1,
	}

	h.renderer.RenderPage(w, "people_list", data)
}

func (h *PersonHandler) ShowPerson(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	person, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	contacts, err := h.contactRepo.GetByPersonID(id)
	if err != nil {
		http.Error(w, "Failed to fetch contacts", http.StatusInternalServerError)
		return
	}

	contactTypes, err := h.contactRepo.GetContactTypes()
	if err != nil {
		http.Error(w, "Failed to fetch contact types", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title        string
		Person       *models.Person
		Contacts     []models.Contact
		ContactTypes []models.ContactType
	}{
		Title:        "Person Details",
		Person:       person,
		Contacts:     contacts,
		ContactTypes: contactTypes,
	}

	h.renderer.RenderPage(w, "person_detail", data)
}

func (h *PersonHandler) NewPersonForm(w http.ResponseWriter, r *http.Request) {
	data := struct {
		Title  string
		Person *models.Person
		Form   *models.PersonForm
	}{
		Title:  "Add New Person",
		Person: nil,
		Form:   &models.PersonForm{},
	}

	h.renderer.RenderPage(w, "person_form", data)
}

func (h *PersonHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	form := &models.PersonForm{
		FirstName:  r.FormValue("firstName"),
		SecondName: r.FormValue("secondName"),
		MiddleName: r.FormValue("middleName"),
		Birthdate:  r.FormValue("birthdate"),
	}

	person := form.ToPerson()
	if err := h.repo.Create(person); err != nil {
		http.Error(w, "Failed to create person", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/people/%d", person.ID))
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/people/%d", person.ID), http.StatusSeeOther)
}

func (h *PersonHandler) EditPersonForm(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	person, err := h.repo.GetByID(id)
	if err != nil {
		http.Error(w, "Person not found", http.StatusNotFound)
		return
	}

	data := struct {
		Title  string
		Person *models.Person
		Form   *models.PersonForm
	}{
		Title:  "Edit Person",
		Person: person,
		Form:   person.ToForm(),
	}

	h.renderer.RenderPage(w, "person_form", data)
}

func (h *PersonHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	form := &models.PersonForm{
		FirstName:  r.FormValue("firstName"),
		SecondName: r.FormValue("secondName"),
		MiddleName: r.FormValue("middleName"),
		Birthdate:  r.FormValue("birthdate"),
	}

	person := form.ToPerson()
	person.ID = id

	if err := h.repo.Update(person); err != nil {
		http.Error(w, "Failed to update person", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/people/%d", person.ID))
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/people/%d", person.ID), http.StatusSeeOther)
}

func (h *PersonHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid person ID", http.StatusBadRequest)
		return
	}

	if err := h.repo.Delete(id); err != nil {
		http.Error(w, "Failed to delete person", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("HX-Request") == "true" {
		w.Header().Set("HX-Redirect", "/people")
		w.WriteHeader(http.StatusOK)
		return
	}

	http.Redirect(w, r, "/people", http.StatusSeeOther)
}
