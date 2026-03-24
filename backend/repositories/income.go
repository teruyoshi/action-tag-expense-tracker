package repositories

import (
	"action-tag-expense-tracker/backend/models"

	"gorm.io/gorm"
)

type IncomeRepository struct {
	DB *gorm.DB
}

func (r *IncomeRepository) FindAll(year, month int) ([]models.Income, error) {
	var incomes []models.Income
	err := r.DB.Preload("IncomeCategory").
		Where("YEAR(date) = ? AND MONTH(date) = ?", year, month).
		Find(&incomes).Error
	return incomes, err
}

func (r *IncomeRepository) FindByID(id uint) (*models.Income, error) {
	var income models.Income
	err := r.DB.Preload("IncomeCategory").First(&income, id).Error
	return &income, err
}

func (r *IncomeRepository) Create(income *models.Income) error {
	return r.DB.Create(income).Error
}

func (r *IncomeRepository) Update(income *models.Income) error {
	return r.DB.Save(income).Error
}

func (r *IncomeRepository) Delete(id uint) error {
	return r.DB.Delete(&models.Income{}, id).Error
}
