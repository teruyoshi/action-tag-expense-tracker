package repositories

import (
	"action-tag-expense-tracker/backend/models"

	"gorm.io/gorm"
)

type EventRepository struct {
	DB *gorm.DB
}

func (r *EventRepository) Create(event *models.Event) error {
	return r.DB.Create(event).Error
}
