package handlers

import (
	"encoding/json"
	"net/http"

	"action-tag-expense-tracker/backend/repositories"
	"action-tag-expense-tracker/backend/services"
)

type BalanceHandler struct {
	Repo    repositories.BalanceRepo
	Service *services.BalanceService
}

func (h *BalanceHandler) Get(w http.ResponseWriter, r *http.Request) {
	balance, err := h.Repo.Get()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, balance)
}

func (h *BalanceHandler) Update(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Amount int `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	balance, err := h.Service.Update(input.Amount)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, balance)
}
