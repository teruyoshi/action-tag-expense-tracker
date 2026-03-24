package repositories

import "action-tag-expense-tracker/backend/models"

type ActionTagRepo interface {
	FindAll() ([]models.ActionTag, error)
	Create(tag *models.ActionTag) error
	Update(tag *models.ActionTag) error
	Delete(id uint) error
	FindByID(id uint) (*models.ActionTag, error)
	FindOrCreateByName(name string) (*models.ActionTag, error)
}

type EventRepo interface {
	Create(event *models.Event) error
}

type ExpenseRepo interface {
	Create(expense *models.Expense) error
	Update(expense *models.Expense) error
}

type SummaryRepo interface {
	MonthTotal(year, month int) (int, error)
	TagMonthTotals(year, month int) ([]TagSummary, error)
	TagMonthTotalsWithDiff(year, month int) ([]TagSummaryWithDiff, error)
	TagExpenseDetails(year, month, tagID int) ([]TagExpenseDetail, error)
}

type BalanceRepo interface {
	Get() (*models.Balance, error)
	Update(amount int) (*models.Balance, error)
	Subtract(amount int) error
	Add(amount int) error
}

type IncomeCategoryRepo interface {
	FindAll() ([]models.IncomeCategory, error)
	Create(category *models.IncomeCategory) error
	Update(category *models.IncomeCategory) error
	Delete(id uint) error
	FindByID(id uint) (*models.IncomeCategory, error)
}

type IncomeRepo interface {
	FindAll(year, month int) ([]models.Income, error)
	FindByID(id uint) (*models.Income, error)
	Create(income *models.Income) error
	Update(income *models.Income) error
	Delete(id uint) error
}
