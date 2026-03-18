package handlers

import (
	"net/http"
	"strconv"

	"action-tag-expense-tracker/backend/repositories"
)

type SummaryHandler struct {
	Repo repositories.SummaryRepo
}

func (h *SummaryHandler) MonthTotal(w http.ResponseWriter, r *http.Request) {
	year, month, err := parseYearMonth(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "year and month query params required")
		return
	}
	total, err := h.Repo.MonthTotal(year, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, map[string]int{"total": total})
}

func (h *SummaryHandler) TagMonthTotals(w http.ResponseWriter, r *http.Request) {
	year, month, err := parseYearMonth(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "year and month query params required")
		return
	}
	results, err := h.Repo.TagMonthTotals(year, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, results)
}

func (h *SummaryHandler) TagMonthTotalsWithDiff(w http.ResponseWriter, r *http.Request) {
	year, month, err := parseYearMonth(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "year and month query params required")
		return
	}
	results, err := h.Repo.TagMonthTotalsWithDiff(year, month)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, results)
}

func (h *SummaryHandler) TagExpenseDetails(w http.ResponseWriter, r *http.Request) {
	year, month, err := parseYearMonth(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "year and month query params required")
		return
	}
	tagID, err := strconv.Atoi(r.URL.Query().Get("tag_id"))
	if err != nil {
		writeError(w, http.StatusBadRequest, "tag_id query param required")
		return
	}
	results, err := h.Repo.TagExpenseDetails(year, month, tagID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	writeJSON(w, results)
}
