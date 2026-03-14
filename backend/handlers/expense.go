package handlers

import (
	"encoding/json"
	"net/http"

	"action-tag-expense-tracker/backend/models"
	"action-tag-expense-tracker/backend/repositories"
)

type ExpenseHandler struct {
	Repo *repositories.ExpenseRepository
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if expense.Amount <= 0 || expense.EventID == 0 {
		http.Error(w, "event_id and positive amount are required", http.StatusBadRequest)
		return
	}
	if err := h.Repo.Create(&expense); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, expense)
}
