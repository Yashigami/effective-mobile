package handler

import (
	"effective-mobail/internal/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// @Summary Удалить человека
// @Description Удаляет человека по ID
// @Param id path int true "ID человека"
// @Success 200 {string} string "Успешно удалено"
// @Failure 404 {string} string "Человек не найден"
// @Failure 500 {string} string "Ошибка при удалении"
// @Router /people/{id} [delete]
func (h *PeopleHandler) DeletePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)                 // Получаю путь параметры
	id, err := strconv.Atoi(vars["id"]) // Преобразуем ID из строки в int
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	// Удаление
	if err := h.DB.Delete(&model.Person{}, id).Error; err != nil {
		http.Error(w, "Ошибка удаления", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Удаление успешно"))
}

// @Summary Обновить данные человека
// @Description Обновляет информацию о человеке по ID
// @Accept json
// @Produce json
// @Param id path int true "ID человека"
// @Param person body model.Person true "Обновлённые данные"
// @Success 200 {object} model.Person
// @Failure 400 {string} string "Неверные данные"
// @Failure 404 {string} string "Человек не найден"
// @Failure 500 {string} string "Внутренняя ошибка сервера"
// @Router /people/{id} [put]
func (h *PeopleHandler) UpdatePerson(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Неверный ID", http.StatusBadRequest)
		return
	}

	var input model.Person
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Невалидный JSON", http.StatusBadRequest)
		return
	}

	var person model.Person
	if err := h.DB.First(&person, id).Error; err != nil {
		http.Error(w, "Человек не найден", http.StatusNotFound)
		return
	}

	// Обновляем нужные поля
	person.Name = input.Name
	person.Surname = input.Surname
	person.Patronymic = input.Patronymic

	// Обновление в БД
	if err := h.DB.Save(&person).Error; err != nil {
		http.Error(w, "Ошибка  при обновлении", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(person)
}
