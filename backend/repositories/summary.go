package repositories

import (
	"gorm.io/gorm"
)

type SummaryRepository struct {
	DB *gorm.DB
}

type TagSummary struct {
	Tag   string `json:"tag"`
	Total int    `json:"total"`
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
		SELECT at.name AS tag, COALESCE(SUM(e.amount), 0) AS total
		FROM expenses e
		JOIN events ev ON e.event_id = ev.id
		JOIN action_tags at ON ev.action_tag_id = at.id
		WHERE YEAR(ev.date) = ? AND MONTH(ev.date) = ?
		GROUP BY at.id, at.name
		ORDER BY total DESC
	`, year, month).Scan(&results).Error
	return results, err
}
