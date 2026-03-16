package repositories

import (
	"sort"

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

type TagSummaryWithDiff struct {
	TagID     uint   `json:"tag_id"`
	Tag       string `json:"tag"`
	Total     int    `json:"total"`
	PrevTotal int    `json:"prev_total"`
	Diff      int    `json:"diff"`
}

func (r *SummaryRepository) TagMonthTotalsWithDiff(year, month int) ([]TagSummaryWithDiff, error) {
	current, err := r.TagMonthTotals(year, month)
	if err != nil {
		return nil, err
	}

	prevYear, prevMonth := year, month-1
	if prevMonth < 1 {
		prevYear--
		prevMonth = 12
	}

	prev, err := r.TagMonthTotals(prevYear, prevMonth)
	if err != nil {
		return nil, err
	}

	prevMap := make(map[uint]TagSummary)
	for _, p := range prev {
		prevMap[p.TagID] = p
	}

	seen := make(map[uint]bool)
	var results []TagSummaryWithDiff

	for _, c := range current {
		seen[c.TagID] = true
		p := prevMap[c.TagID]
		results = append(results, TagSummaryWithDiff{
			TagID:     c.TagID,
			Tag:       c.Tag,
			Total:     c.Total,
			PrevTotal: p.Total,
			Diff:      c.Total - p.Total,
		})
	}

	for _, p := range prev {
		if seen[p.TagID] {
			continue
		}
		results = append(results, TagSummaryWithDiff{
			TagID:     p.TagID,
			Tag:       p.Tag,
			Total:     0,
			PrevTotal: p.Total,
			Diff:      -p.Total,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Total > results[j].Total
	})

	return results, nil
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
