package repositories

import "action-tag-expense-tracker/backend/models"

type ActionTagRepo interface {
	FindAll() ([]models.ActionTag, error)
	Create(tag *models.ActionTag) error
	Update(tag *models.ActionTag) error
	Delete(id uint) error
	FindByID(id uint) (*models.ActionTag, error)
}

type EventRepo interface {
	Create(event *models.Event) error
}

type ExpenseRepo interface {
	Create(expense *models.Expense) error
}

type SummaryRepo interface {
	MonthTotal(year, month int) (int, error)
	TagMonthTotals(year, month int) ([]TagSummary, error)
}

type BalanceRepo interface {
	Get() (*models.Balance, error)
	Update(amount int) (*models.Balance, error)
	Subtract(amount int) error
}
