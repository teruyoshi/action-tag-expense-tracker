package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"action-tag-expense-tracker/backend/models"
	"action-tag-expense-tracker/backend/repositories"

	"github.com/go-chi/chi/v5"
)

type ExpenseHandler struct {
	Repo        repositories.ExpenseRepo
	BalanceRepo repositories.BalanceRepo
}

func (h *ExpenseHandler) Create(w http.ResponseWriter, r *http.Request) {
	var expense models.Expense
	if err := json.NewDecoder(r.Body).Decode(&expense); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if expense.Amount <= 0 || expense.EventID == 0 {
		writeError(w, http.StatusBadRequest, "event_id and positive amount are required")
		return
	}
	if err := h.Repo.Create(&expense); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if h.BalanceRepo != nil {
		if err := h.BalanceRepo.Subtract(expense.Amount); err != nil {
			writeError(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, expense)
}

func (h *ExpenseHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var input struct {
		Item   string `json:"item"`
		Amount int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "positive amount is required")
		return
	}
	expense := &models.Expense{
		ID:     uint(id),
		Item:   input.Item,
		Amount: input.Amount,
	}
	if err := h.Repo.Update(expense); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, expense)
}
