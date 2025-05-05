package handler

import (
	"effective-mobail/internal/model"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

// Эндпоинт DELETE /people/{id}
// DeletePerson удаляет человека по ID

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

//PUT /people/{id}

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
