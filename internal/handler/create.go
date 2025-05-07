package handler

import (
	"effective-mobail/internal/model"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
)

// создаем валидатор
var validate = validator.New()

// CreatePersonHandler обрабатывает создание человека
// @Summary Создание человека
// @Accept json
// @Produce json
// @Param person body model.Person true "Данные человека"
// @Success 200 {object} model.Person
// @Failure 400 {string} string "Неверные данные"
// @Router /people [post]
func (h *PeopleHandler) CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person model.Person // создаём переменную, куда положим JSON-данные

	// Декодируем JSON из тела запроса в структуру person
	if err := json.NewDecoder(r.Body).Decode(&person); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}


	// Валидируем поля структуры (используем валидатор)
	if err := validate.Struct(person); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Обогащаем человека (возраст, пол, национальность)
	person = EnrichPerson(person)

	// Сохраняем в БД, обращаясь к полю h.DB
	if err := h.DB.Create(&person).Error; err != nil {
		http.Error(w, "Сохранение не удалось", http.StatusInternalServerError)
		return
	}

	// Отправляем успешный JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}
