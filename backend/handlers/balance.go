package handlers

import (
	"encoding/json"
	"net/http"
	"time"

	"action-tag-expense-tracker/backend/models"
	"action-tag-expense-tracker/backend/repositories"
)

type BalanceHandler struct {
	Repo          repositories.BalanceRepo
	ActionTagRepo repositories.ActionTagRepo
	EventRepo     repositories.EventRepo
	ExpenseRepo   repositories.ExpenseRepo
}

func (h *BalanceHandler) Get(w http.ResponseWriter, r *http.Request) {
	balance, err := h.Repo.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, balance)
}

func (h *BalanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	old, err := h.Repo.Get()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	balance, err := h.Repo.Update(input.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	diff := old.Amount - input.Amount
	if diff > 0 && h.ActionTagRepo != nil && h.EventRepo != nil && h.ExpenseRepo != nil {
		tag, err := h.ActionTagRepo.FindOrCreateByName("その他")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		event := &models.Event{
			Date:        time.Now(),
			ActionTagID: tag.ID,
		}
		if err := h.EventRepo.Create(event); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		expense := &models.Expense{
			EventID: event.ID,
			Item:    "",
			Amount:  diff,
		}
		if err := h.ExpenseRepo.Create(expense); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}

	writeJSON(w, balance)
}
