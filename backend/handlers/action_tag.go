package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"action-tag-expense-tracker/backend/models"
	"action-tag-expense-tracker/backend/repositories"

	"github.com/go-chi/chi/v5"
)

type ActionTagHandler struct {
	Repo repositories.ActionTagRepo
}

func (h *ActionTagHandler) List(w http.ResponseWriter, r *http.Request) {
	tags, err := h.Repo.FindAll()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, tags)
}

func (h *ActionTagHandler) Create(w http.ResponseWriter, r *http.Request) {
	var tag models.ActionTag
	if err := json.NewDecoder(r.Body).Decode(&tag); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if tag.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	if err := h.Repo.Create(&tag); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, tag)
}

func (h *ActionTagHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	tag, err := h.Repo.FindByID(uint(id))
	if err != nil {
		writeError(w, http.StatusNotFound, "not found")
		return
	}
	var input struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}
	tag.Name = input.Name
	if err := h.Repo.Update(tag); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, tag)
}

func (h *ActionTagHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.Repo.Delete(uint(id)); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
