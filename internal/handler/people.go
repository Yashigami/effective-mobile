package handler

import (
	"effective-mobail/internal/model"
	"encoding/json"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

// Обертка с доступом к БД
type PeopleHandler struct {
	DB *gorm.DB
}

// NewPeopleHandler создает новый хендлер
func NewPeopleHandler(db *gorm.DB) *PeopleHandler {
	return &PeopleHandler{DB: db}
}

// CreatePerson обрабатывает POST-запрос на создание человека
func (h *PeopleHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var input model.Person // Создаем пустую структуру

	// Читаем тело запроса
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Ошибка чтения запроса", http.StatusBadRequest)
		return
	}

	// Парсим JSON в структуру
	if err := json.Unmarshal(body, &input); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	// Проверка: имя обязательно
	if input.Name == "" {
		http.Error(w, "Имя обязательно", http.StatusBadRequest)
		return
	}
	log.Printf("[INFO] Получен запрос на создание: %+v\n", input)

	// Обогащение через API
	enrichedPerson := EnrichPerson(input)

	// Сохранаю в базу
	result := h.DB.Create(&enrichedPerson)
	if result.Error != nil {
		http.Error(w, "Ошибка при сохранение в БД", http.StatusInternalServerError)
		return
	}

	// Возвращаемый успешный ответ
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(enrichedPerson)
}
