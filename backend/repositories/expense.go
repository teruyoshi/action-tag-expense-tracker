package repositories

import (
	"action-tag-expense-tracker/backend/models"

	"gorm.io/gorm"
)

type ExpenseRepository struct {
	DB *gorm.DB
}

func (r *ExpenseRepository) Create(expense *models.Expense) error {
	return r.DB.Create(expense).Error
}
