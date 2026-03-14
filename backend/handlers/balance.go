package handlers

import (
	"encoding/json"
	"net/http"

	"action-tag-expense-tracker/backend/repositories"
)

type BalanceHandler struct {
	Repo repositories.BalanceRepo
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
	balance, err := h.Repo.Update(input.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, balance)
}
