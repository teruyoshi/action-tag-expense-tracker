package services

import (
	"testing"
	"time"

	"action-tag-expense-tracker/backend/models"
)

// --- mocks ---

type mockIncomeRepo struct {
	income  *models.Income
	incomes []models.Income
	created *models.Income
	updated *models.Income
	deleted uint
	err     error
}

func (m *mockIncomeRepo) FindAll(year, month int) ([]models.Income, error) {
	return m.incomes, m.err
}
func (m *mockIncomeRepo) FindByID(id uint) (*models.Income, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.income, nil
}
func (m *mockIncomeRepo) Create(income *models.Income) error {
	if m.err != nil {
		return m.err
	}
	income.ID = 1
	m.created = income
	return nil
}
func (m *mockIncomeRepo) Update(income *models.Income) error {
	if m.err != nil {
		return m.err
	}
	m.updated = income
	return nil
}
func (m *mockIncomeRepo) Delete(id uint) error {
	if m.err != nil {
		return m.err
	}
	m.deleted = id
	return nil
}

// --- tests ---

func TestIncomeService_Create(t *testing.T) {
	tests := []struct {
		name        string
		incomeRepo  *mockIncomeRepo
		balanceRepo *mockBalanceRepo
		income      *models.Income
		wantErr     bool
		wantBalance int
	}{
		{
			name:        "creates income and adds to balance",
			incomeRepo:  &mockIncomeRepo{},
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			income: &models.Income{
				IncomeCategoryID: 1,
				Date:             time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC),
				Description:      "給与",
				Amount:           300000,
			},
			wantBalance: 400000,
		},
		{
			name:        "returns error on income repo failure",
			incomeRepo:  &mockIncomeRepo{err: errDB},
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			income: &models.Income{
				Amount: 300000,
			},
			wantErr: true,
		},
		{
			name:        "returns error on balance repo failure",
			incomeRepo:  &mockIncomeRepo{},
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}, err: errDB},
			income: &models.Income{
				Amount: 300000,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Balance repo error should only apply to Add, not Create
			balanceErr := tt.balanceRepo.err
			if tt.name == "returns error on balance repo failure" {
				tt.balanceRepo.err = nil
			}

			s := &IncomeService{
				IncomeRepo:  tt.incomeRepo,
				BalanceRepo: tt.balanceRepo,
			}

			if tt.name == "returns error on balance repo failure" {
				tt.balanceRepo.err = balanceErr
			}

			err := s.Create(tt.income)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.balanceRepo.balance.Amount != tt.wantBalance {
				t.Errorf("balance = %d, want %d", tt.balanceRepo.balance.Amount, tt.wantBalance)
			}
		})
	}
}

func TestIncomeService_Update(t *testing.T) {
	tests := []struct {
		name        string
		incomeRepo  *mockIncomeRepo
		balanceRepo *mockBalanceRepo
		income      *models.Income
		wantErr     bool
		wantBalance int
	}{
		{
			name: "increase amount adds to balance",
			incomeRepo: &mockIncomeRepo{
				income: &models.Income{ID: 1, Amount: 200000},
			},
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 300000}},
			income:      &models.Income{ID: 1, IncomeCategoryID: 1, Date: time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC), Amount: 250000},
			wantBalance: 350000,
		},
		{
			name: "decrease amount subtracts from balance",
			incomeRepo: &mockIncomeRepo{
				income: &models.Income{ID: 1, Amount: 300000},
			},
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 400000}},
			income:      &models.Income{ID: 1, IncomeCategoryID: 1, Date: time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC), Amount: 250000},
			wantBalance: 350000,
		},
		{
			name: "same amount no balance change",
			incomeRepo: &mockIncomeRepo{
				income: &models.Income{ID: 1, Amount: 300000},
			},
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 400000}},
			income:      &models.Income{ID: 1, IncomeCategoryID: 1, Date: time.Date(2026, 3, 25, 0, 0, 0, 0, time.UTC), Amount: 300000},
			wantBalance: 400000,
		},
		{
			name:        "returns error on FindByID failure",
			incomeRepo:  &mockIncomeRepo{err: errDB},
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			income:      &models.Income{ID: 1, Amount: 250000},
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IncomeService{
				IncomeRepo:  tt.incomeRepo,
				BalanceRepo: tt.balanceRepo,
			}

			err := s.Update(tt.income)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.balanceRepo.balance.Amount != tt.wantBalance {
				t.Errorf("balance = %d, want %d", tt.balanceRepo.balance.Amount, tt.wantBalance)
			}
		})
	}
}

func TestIncomeService_Delete(t *testing.T) {
	tests := []struct {
		name        string
		incomeRepo  *mockIncomeRepo
		balanceRepo *mockBalanceRepo
		id          uint
		wantErr     bool
		wantBalance int
	}{
		{
			name: "deletes income and subtracts from balance",
			incomeRepo: &mockIncomeRepo{
				income: &models.Income{ID: 1, Amount: 300000},
			},
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 400000}},
			id:          1,
			wantBalance: 100000,
		},
		{
			name:        "returns error on FindByID failure",
			incomeRepo:  &mockIncomeRepo{err: errDB},
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			id:          1,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &IncomeService{
				IncomeRepo:  tt.incomeRepo,
				BalanceRepo: tt.balanceRepo,
			}

			err := s.Delete(tt.id)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !tt.wantErr && tt.balanceRepo.balance.Amount != tt.wantBalance {
				t.Errorf("balance = %d, want %d", tt.balanceRepo.balance.Amount, tt.wantBalance)
			}
		})
	}
}
