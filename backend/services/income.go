package services

import (
	"action-tag-expense-tracker/backend/models"
	"action-tag-expense-tracker/backend/repositories"
)

type IncomeService struct {
	IncomeRepo  repositories.IncomeRepo
	BalanceRepo repositories.BalanceRepo
}

func (s *IncomeService) Create(income *models.Income) error {
	if err := s.IncomeRepo.Create(income); err != nil {
		return err
	}
	return s.BalanceRepo.Add(income.Amount)
}

func (s *IncomeService) Update(income *models.Income) error {
	old, err := s.IncomeRepo.FindByID(income.ID)
	if err != nil {
		return err
	}

	diff := income.Amount - old.Amount
	if err := s.IncomeRepo.Update(income); err != nil {
		return err
	}

	if diff > 0 {
		return s.BalanceRepo.Add(diff)
	} else if diff < 0 {
		return s.BalanceRepo.Subtract(-diff)
	}
	return nil
}

func (s *IncomeService) Delete(id uint) error {
	income, err := s.IncomeRepo.FindByID(id)
	if err != nil {
		return err
	}

	if err := s.IncomeRepo.Delete(id); err != nil {
		return err
	}
	return s.BalanceRepo.Subtract(income.Amount)
}
