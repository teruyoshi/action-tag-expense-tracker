package models

import "time"

type ActionTag struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"not null"`
}

type Event struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Date        time.Time `json:"date" gorm:"type:date;not null"`
	ActionTagID uint      `json:"action_tag_id" gorm:"not null"`
	ActionTag   ActionTag `json:"action_tag,omitempty" gorm:"foreignKey:ActionTagID"`
	Expenses    []Expense `json:"expenses,omitempty" gorm:"foreignKey:EventID"`
}

type Expense struct {
	ID      uint   `json:"id" gorm:"primaryKey"`
	EventID uint   `json:"event_id" gorm:"not null"`
	Item    string `json:"item" gorm:"not null"`
	Amount  int    `json:"amount" gorm:"not null"`
}
