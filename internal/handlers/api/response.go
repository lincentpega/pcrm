package api

import (
	"encoding/json"
	"net/http"
)

type PaginatedResponse[T any] struct {
	Data        []T  `json:"data"`
	CurrentPage int  `json:"currentPage"`
	TotalPages  int  `json:"totalPages"`
	TotalCount  int  `json:"totalCount"`
	HasNext     bool `json:"hasNext"`
	HasPrev     bool `json:"hasPrev"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func WriteJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func WriteSuccess[T any](w http.ResponseWriter, data T) {
	WriteJSON(w, http.StatusOK, data)
}

func WriteCreated[T any](w http.ResponseWriter, data T) {
	WriteJSON(w, http.StatusCreated, data)
}

func WritePaginated[T any](w http.ResponseWriter, data []T, currentPage, totalPages, totalCount int) {
	response := PaginatedResponse[T]{
		Data:        data,
		CurrentPage: currentPage,
		TotalPages:  totalPages,
		TotalCount:  totalCount,
		HasNext:     currentPage < totalPages,
		HasPrev:     currentPage > 1,
	}
	WriteJSON(w, http.StatusOK, response)
}

func WriteError(w http.ResponseWriter, status int, err string) {
	response := ErrorResponse{
		Error: err,
	}
	WriteJSON(w, status, response)
}

func WriteBadRequest(w http.ResponseWriter, err string) {
	WriteError(w, http.StatusBadRequest, err)
}

func WriteNotFound(w http.ResponseWriter, err string) {
	WriteError(w, http.StatusNotFound, err)
}

func WriteInternalError(w http.ResponseWriter, err string) {
	WriteError(w, http.StatusInternalServerError, err)
}