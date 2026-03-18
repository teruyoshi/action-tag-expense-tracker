package services

import (
	"time"

	"action-tag-expense-tracker/backend/models"
	"action-tag-expense-tracker/backend/repositories"
)

type BalanceService struct {
	BalanceRepo   repositories.BalanceRepo
	ActionTagRepo repositories.ActionTagRepo
	EventRepo     repositories.EventRepo
	ExpenseRepo   repositories.ExpenseRepo
}

func (s *BalanceService) Update(newAmount int) (*models.Balance, error) {
	old, err := s.BalanceRepo.Get()
	if err != nil {
		return nil, err
	}

	diff := old.Amount - newAmount
	if diff > 0 {
		if err := s.createExpense(diff); err != nil {
			return nil, err
		}
	}

	return s.BalanceRepo.Update(newAmount)
}

func (s *BalanceService) createExpense(amount int) error {
	tag, err := s.ActionTagRepo.FindOrCreateByName("その他")
	if err != nil {
		return err
	}

	event := &models.Event{
		Date:        time.Now(),
		ActionTagID: tag.ID,
	}
	if err := s.EventRepo.Create(event); err != nil {
		return err
	}

	expense := &models.Expense{
		EventID: event.ID,
		Item:    "",
		Amount:  amount,
	}
	return s.ExpenseRepo.Create(expense)
}
