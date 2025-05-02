package model

import "gorm.io/gorm"

// Person описывает данные о человеке
type Person struct {
	gorm.Model // встроенные поля с ID, CreateAt, UpdateAt, DeleteAt

	Name        string  `json:"name"`
	Surname     string  `json:"surname"`
	Patronymic  *string `json:"patronymic"`
	Gender      *string `json:"gender"`
	Age         *int    `json:"age"`
	Nationality *string `json:"nationality"`
}
