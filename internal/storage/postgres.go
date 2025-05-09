package storage

import (
	"effective-mobail/internal/config"
	"effective-mobail/internal/model"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type Storage struct {
	DB *gorm.DB
}

func (s *Storage) SavePerson(p *model.Person) error {
	return s.DB.Create(p).Error // 👈 теперь мы сохраняем через GORM
}

func InitPostgres(cfg *config.Config) *gorm.DB {
	// Строка подключения к PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", cfg.DBHost, cfg.DBUser, cfg.DBPass, cfg.DBName, cfg.DBPort)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("неудалось подключиться к БД: %v", err)
	}

	// Автоматическая миграция таблицы Person
	err = db.AutoMigrate(&model.Person{})
	if err != nil {
		log.Fatalf("миграция не удалась: %v", err)
	}
	log.Println("Успешное подключение к PostgreSQL и миграция")
	return db
}
