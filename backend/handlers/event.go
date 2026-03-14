package handlers

import (
	"encoding/json"
	"net/http"

	"action-tag-expense-tracker/backend/models"
	"action-tag-expense-tracker/backend/repositories"
)

type EventHandler struct {
	Repo repositories.EventRepo
}

type CreateEventRequest struct {
	Date        string `json:"date"`
	ActionTagID uint   `json:"action_tag_id"`
}

func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	date, err := parseDate(req.Date)
	if err != nil {
		http.Error(w, "invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	event := models.Event{
		Date:        date,
		ActionTagID: req.ActionTagID,
	}
	if err := h.Repo.Create(&event); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, event)
}
