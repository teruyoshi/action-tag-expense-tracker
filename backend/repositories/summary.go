package repositories

import (
	"gorm.io/gorm"
)

type SummaryRepository struct {
	DB *gorm.DB
}

type TagSummary struct {
	TagID uint   `json:"tag_id"`
	Tag   string `json:"tag"`
	Total int    `json:"total"`
}

type TagExpenseDetail struct {
	ID     uint   `json:"id"`
	Date   string `json:"date"`
	Item   string `json:"item"`
	Amount int    `json:"amount"`
}

func (r *SummaryRepository) MonthTotal(year, month int) (int, error) {
	var total int
	err := r.DB.Raw(`
		SELECT COALESCE(SUM(e.amount), 0)
		FROM expenses e
		JOIN events ev ON e.event_id = ev.id
		WHERE YEAR(ev.date) = ? AND MONTH(ev.date) = ?
	`, year, month).Scan(&total).Error
	return total, err
}

func (r *SummaryRepository) TagMonthTotals(year, month int) ([]TagSummary, error) {
	var results []TagSummary
	err := r.DB.Raw(`
		SELECT at.id AS tag_id, at.name AS tag, COALESCE(SUM(e.amount), 0) AS total
		FROM expenses e
		JOIN events ev ON e.event_id = ev.id
		JOIN action_tags at ON ev.action_tag_id = at.id
		WHERE YEAR(ev.date) = ? AND MONTH(ev.date) = ?
		GROUP BY at.id, at.name
		ORDER BY total DESC
	`, year, month).Scan(&results).Error
	return results, err
}

func (r *SummaryRepository) TagExpenseDetails(year, month, tagID int) ([]TagExpenseDetail, error) {
	var results []TagExpenseDetail
	err := r.DB.Raw(`
		SELECT e.id, DATE_FORMAT(ev.date, '%Y-%m-%d') AS date, e.item, e.amount
		FROM expenses e
		JOIN events ev ON e.event_id = ev.id
		WHERE YEAR(ev.date) = ? AND MONTH(ev.date) = ? AND ev.action_tag_id = ?
		ORDER BY ev.date DESC, e.id DESC
	`, year, month, tagID).Scan(&results).Error
	return results, err
}
