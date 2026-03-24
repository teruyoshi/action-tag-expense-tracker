package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"action-tag-expense-tracker/backend/models"
	"action-tag-expense-tracker/backend/repositories"
	"action-tag-expense-tracker/backend/services"

	"github.com/go-chi/chi/v5"
)

type IncomeHandler struct {
	Repo    repositories.IncomeRepo
	Service *services.IncomeService
}

func (h *IncomeHandler) List(w http.ResponseWriter, r *http.Request) {
	year, month, err := parseYearMonth(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "year and month are required")
		return
	}
	incomes, err := h.Repo.FindAll(year, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, incomes)
}

func (h *IncomeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var input struct {
		IncomeCategoryID uint   `json:"income_category_id"`
		Date             string `json:"date"`
		Description      string `json:"description"`
		Amount           int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.Amount <= 0 || input.IncomeCategoryID == 0 || input.Date == "" {
		writeError(w, http.StatusBadRequest, "income_category_id, date and positive amount are required")
		return
	}
	date, err := parseDate(input.Date)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format (expected YYYY-MM-DD)")
		return
	}

	income := &models.Income{
		IncomeCategoryID: input.IncomeCategoryID,
		Date:             date,
		Description:      input.Description,
		Amount:           input.Amount,
	}
	if err := h.Service.Create(income); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusCreated)
	writeJSON(w, income)
}

func (h *IncomeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	var input struct {
		IncomeCategoryID uint   `json:"income_category_id"`
		Date             string `json:"date"`
		Description      string `json:"description"`
		Amount           int    `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if input.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "positive amount is required")
		return
	}
	date, err := parseDate(input.Date)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid date format (expected YYYY-MM-DD)")
		return
	}

	income := &models.Income{
		ID:               uint(id),
		IncomeCategoryID: input.IncomeCategoryID,
		Date:             date,
		Description:      input.Description,
		Amount:           input.Amount,
	}
	if err := h.Service.Update(income); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, income)
}

func (h *IncomeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil || id < 1 {
		writeError(w, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.Service.Delete(uint(id)); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
