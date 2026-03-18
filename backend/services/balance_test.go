package services

import (
	"errors"
	"testing"

	"action-tag-expense-tracker/backend/models"
)

var errDB = errors.New("db error")

// --- mocks ---

type mockBalanceRepo struct {
	balance *models.Balance
	err     error
}

func (m *mockBalanceRepo) Get() (*models.Balance, error) {
	return m.balance, m.err
}
func (m *mockBalanceRepo) Update(amount int) (*models.Balance, error) {
	if m.err != nil {
		return nil, m.err
	}
	return &models.Balance{ID: 1, Amount: amount}, nil
}
func (m *mockBalanceRepo) Subtract(amount int) error {
	if m.err != nil {
		return m.err
	}
	m.balance.Amount -= amount
	return nil
}

type mockActionTagRepo struct {
	tags []models.ActionTag
	err  error
}

func (m *mockActionTagRepo) FindAll() ([]models.ActionTag, error) { return m.tags, m.err }
func (m *mockActionTagRepo) Create(tag *models.ActionTag) error   { return m.err }
func (m *mockActionTagRepo) Update(tag *models.ActionTag) error   { return m.err }
func (m *mockActionTagRepo) Delete(id uint) error                 { return m.err }
func (m *mockActionTagRepo) FindByID(id uint) (*models.ActionTag, error) {
	return nil, m.err
}
func (m *mockActionTagRepo) FindOrCreateByName(name string) (*models.ActionTag, error) {
	if m.err != nil {
		return nil, m.err
	}
	for i := range m.tags {
		if m.tags[i].Name == name {
			return &m.tags[i], nil
		}
	}
	tag := models.ActionTag{ID: 1, Name: name}
	return &tag, nil
}

type mockEventRepo struct {
	created *models.Event
	err     error
}

func (m *mockEventRepo) Create(event *models.Event) error {
	if m.err != nil {
		return m.err
	}
	event.ID = 1
	m.created = event
	return nil
}

type mockExpenseRepo struct {
	created *models.Expense
	err     error
}

func (m *mockExpenseRepo) Create(expense *models.Expense) error {
	if m.err != nil {
		return m.err
	}
	expense.ID = 1
	m.created = expense
	return nil
}
func (m *mockExpenseRepo) Update(expense *models.Expense) error { return m.err }

// --- tests ---

func TestBalanceService_Update(t *testing.T) {
	tests := []struct {
		name           string
		balanceRepo    *mockBalanceRepo
		tagRepo        *mockActionTagRepo
		eventRepo      *mockEventRepo
		expenseRepo    *mockExpenseRepo
		newAmount      int
		wantAmount     int
		wantErr        bool
		wantEvent      bool
		wantExpense    bool
		wantExpenseAmt int
	}{
		{
			name:        "increase - no expense created",
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			tagRepo:     &mockActionTagRepo{},
			eventRepo:   &mockEventRepo{},
			expenseRepo: &mockExpenseRepo{},
			newAmount:   150000,
			wantAmount:  150000,
			wantEvent:   false,
			wantExpense: false,
		},
		{
			name:           "decrease - expense created",
			balanceRepo:    &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			tagRepo:        &mockActionTagRepo{tags: []models.ActionTag{{ID: 5, Name: "その他"}}},
			eventRepo:      &mockEventRepo{},
			expenseRepo:    &mockExpenseRepo{},
			newAmount:      80000,
			wantAmount:     80000,
			wantEvent:      true,
			wantExpense:    true,
			wantExpenseAmt: 20000,
		},
		{
			name:        "no change - no expense created",
			balanceRepo: &mockBalanceRepo{balance: &models.Balance{ID: 1, Amount: 100000}},
			tagRepo:     &mockActionTagRepo{},
			eventRepo:   &mockEventRepo{},
			expenseRepo: &mockExpenseRepo{},
			newAmount:   100000,
			wantAmount:  100000,
			wantEvent:   false,
			wantExpense: false,
		},
		{
			name:        "returns error on Get failure",
			balanceRepo: &mockBalanceRepo{err: errDB},
			tagRepo:     &mockActionTagRepo{},
			eventRepo:   &mockEventRepo{},
			expenseRepo: &mockExpenseRepo{},
			newAmount:   100000,
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &BalanceService{
				BalanceRepo:   tt.balanceRepo,
				ActionTagRepo: tt.tagRepo,
				EventRepo:     tt.eventRepo,
				ExpenseRepo:   tt.expenseRepo,
			}

			result, err := s.Update(tt.newAmount)
			if (err != nil) != tt.wantErr {
				t.Fatalf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.wantErr {
				return
			}

			if result.Amount != tt.wantAmount {
				t.Errorf("amount = %d, want %d", result.Amount, tt.wantAmount)
			}
			if tt.wantEvent && tt.eventRepo.created == nil {
				t.Error("expected event to be created, but it was not")
			}
			if !tt.wantEvent && tt.eventRepo.created != nil {
				t.Error("expected no event, but one was created")
			}
			if tt.wantExpense {
				if tt.expenseRepo.created == nil {
					t.Error("expected expense to be created, but it was not")
				} else {
					if tt.expenseRepo.created.Amount != tt.wantExpenseAmt {
						t.Errorf("expense amount = %d, want %d", tt.expenseRepo.created.Amount, tt.wantExpenseAmt)
					}
					if tt.expenseRepo.created.Item != "" {
						t.Errorf("expense item = %q, want empty string", tt.expenseRepo.created.Item)
					}
				}
			}
			if !tt.wantExpense && tt.expenseRepo.created != nil {
				t.Error("expected no expense, but one was created")
			}
		})
	}
}
