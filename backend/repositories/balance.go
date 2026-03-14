package repositories

import (
	"action-tag-expense-tracker/backend/models"

	"gorm.io/gorm"
)

type BalanceRepository struct {
	DB *gorm.DB
}

func (r *BalanceRepository) Get() (*models.Balance, error) {
	var balance models.Balance
	err := r.DB.First(&balance).Error
	if err == gorm.ErrRecordNotFound {
		balance = models.Balance{Amount: 0}
		if createErr := r.DB.Create(&balance).Error; createErr != nil {
			return nil, createErr
		}
		return &balance, nil
	}
	return &balance, err
}

func (r *BalanceRepository) Update(amount int) (*models.Balance, error) {
	balance, err := r.Get()
	if err != nil {
		return nil, err
	}
	balance.Amount = amount
	if err := r.DB.Save(balance).Error; err != nil {
		return nil, err
	}
	return balance, nil
}

func (r *BalanceRepository) Subtract(amount int) error {
	balance, err := r.Get()
	if err != nil {
		return err
	}
	balance.Amount -= amount
	return r.DB.Save(balance).Error
}
