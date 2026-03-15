package handlers

import (
	"errors"

	"action-tag-expense-tracker/backend/models"
	"action-tag-expense-tracker/backend/repositories"
)

var errDB = errors.New("db error")

// --- ActionTagRepo mock ---

type mockActionTagRepo struct {
	tags    []models.ActionTag
	created *models.ActionTag
	err     error
}

func (m *mockActionTagRepo) FindAll() ([]models.ActionTag, error) {
	return m.tags, m.err
}
func (m *mockActionTagRepo) Create(tag *models.ActionTag) error {
	if m.err != nil {
		return m.err
	}
	tag.ID = 1
	m.created = tag
	return nil
}
func (m *mockActionTagRepo) Update(tag *models.ActionTag) error {
	return m.err
}
func (m *mockActionTagRepo) Delete(id uint) error {
	return m.err
}
func (m *mockActionTagRepo) FindByID(id uint) (*models.ActionTag, error) {
	if m.err != nil {
		return nil, m.err
	}
	for i := range m.tags {
		if m.tags[i].ID == id {
			return &m.tags[i], nil
		}
	}
	return nil, errors.New("not found")
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
	m.created = &tag
	return &tag, nil
}

// --- EventRepo mock ---

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

// --- ExpenseRepo mock ---

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

func (m *mockExpenseRepo) Update(expense *models.Expense) error {
	return m.err
}

// --- SummaryRepo mock ---

type mockSummaryRepo struct {
	total      int
	tagTotals  []repositories.TagSummary
	tagDetails []repositories.TagExpenseDetail
	err        error
}

func (m *mockSummaryRepo) MonthTotal(year, month int) (int, error) {
	return m.total, m.err
}
func (m *mockSummaryRepo) TagMonthTotals(year, month int) ([]repositories.TagSummary, error) {
	return m.tagTotals, m.err
}
func (m *mockSummaryRepo) TagExpenseDetails(year, month, tagID int) ([]repositories.TagExpenseDetail, error) {
	return m.tagDetails, m.err
}

// --- BalanceRepo mock ---

type mockBalanceRepo struct {
	balance    *models.Balance
	subtracted int
	err        error
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
	m.subtracted = amount
	return nil
}
