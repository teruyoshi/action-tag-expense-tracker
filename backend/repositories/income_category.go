package repositories

import (
	"action-tag-expense-tracker/backend/models"

	"gorm.io/gorm"
)

type IncomeCategoryRepository struct {
	DB *gorm.DB
}

func (r *IncomeCategoryRepository) FindAll() ([]models.IncomeCategory, error) {
	var categories []models.IncomeCategory
	err := r.DB.Find(&categories).Error
	return categories, err
}

func (r *IncomeCategoryRepository) Create(category *models.IncomeCategory) error {
	return r.DB.Create(category).Error
}

func (r *IncomeCategoryRepository) Update(category *models.IncomeCategory) error {
	return r.DB.Save(category).Error
}

func (r *IncomeCategoryRepository) Delete(id uint) error {
	return r.DB.Delete(&models.IncomeCategory{}, id).Error
}

func (r *IncomeCategoryRepository) FindByID(id uint) (*models.IncomeCategory, error) {
	var category models.IncomeCategory
	err := r.DB.First(&category, id).Error
	return &category, err
}
