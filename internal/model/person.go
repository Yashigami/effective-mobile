package model

import (
	"gorm.io/gorm"
	"time"
)

// Person описывает данные о человеке
type Person struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	Name        string         `json:"name" validate:"required"`
	Surname     string         `json:"surname" validate:"required"`
	Patronymic  *string        `json:"patronymic"`
	Age         *int           `json:"age"`
	Gender      *string        `json:"gender"`
	Nationality *string        `json:"nationality"`
}
