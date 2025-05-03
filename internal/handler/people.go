package handler

import (
	"effective-mobail/internal/model"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

// GetPeople обрабатывает GET-запросна получение списка людей
func (h *PeopleHandler) GetPeople(w http.ResponseWriter, r *http.Request) {
	var people []model.Person // сюда загрузим результат
	var total int64

	// Получаем query-параметры
	q := r.URL.Query()

	// Базовый запрос
	db := h.DB.Model(&model.Person{})

	// Фильтрация
	if name := q.Get("name"); name != "" {
		db = db.Where("name ILIKE ?", "%"+name+"%")
	}
	if surname := q.Get("surname"); surname != "" {
		db = db.Where("surname ILIKE ?", "%"+surname+"%")
	}
	if gender := q.Get("gender"); gender != "" {
		db = db.Where("gender = ?", gender)
	}
	if nationality := q.Get("nationality"); nationality != "" {
		db = db.Where("nationality = ?", nationality)
	}

	// Подсчет для общего количества (для фронта)
	db.Count(&total)

	// Пагинация
	limit := 10
	offset := 0
	if val := q.Get("limit"); val != "" {
		fmt.Sscanf(val, "%d", &limit)
	}
	if val := q.Get("offset"); val != "" {
		fmt.Sscanf(val, "%d", &offset)
	}

	// Применяем limit/offset
	db = db.Limit(limit).Offset(offset)

	// Выполняем запрос
	if err := db.Order("name").Find(&people).Error; err != nil {
		http.Error(w, "Ошибка при получение данных", http.StatusInternalServerError)
		return
	}

	// Формируем и отправляем ответ
	response := map[string]interface{}{
		"total":  total,
		"people": people,
	}
	w.Header().Set("Content-Type", "application/json")
	log.Printf("Отправляем ответ: %+v\n", response)

	// отправка JSON
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Ошибка кодирования JSON: %v\n", err)
		http.Error(w, "Ошибка при отправке ответа", http.StatusInternalServerError)
	}
}

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
