package handler

import (
	"effective-mobail/internal/model"
	"effective-mobail/internal/storage"
	"effective-mobail/pkg/enrich"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Структура входного JSON запроса
type enrichRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

// Структура ответа
type enrichResponse struct {
	Name        string `json:"name"`
	Surname     string `json:"surname"`
	Age         int    `json:"age"`
	Gender      string `json:"gender"`
	Nationality string `json:"nationality"`
}

// Хэндлер для маршрута POST /enrich
func EnrichHandler(store *storage.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req enrichRequest

		// 1. Декодируем тело запроса
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "неправильный формат запроса", http.StatusBadRequest)
			return
		}

		// 2. Обогощаем имя через внешние API
		data, err := enrich.Enrich(req.Name)
		if err != nil {
			http.Error(w, "не удалось обогатить данные:"+err.Error(), http.StatusInternalServerError)
			return
		}

		// 3. Сохраняем результат в БД
		person := model.Person{
			Name:        req.Name,
			Surname:     req.Surname,
			Age:         &data.Age,
			Gender:      &data.Gender,
			Nationality: &data.Nationality,
		}
		if err := store.SavePerson(&person); err != nil {
			http.Error(w, "ошибка при сохранение в БД:"+err.Error(), http.StatusInternalServerError)
			return
		}

		// 4. Возвращаем клиенту ответ
		resp := enrichResponse{
			Name:        person.Name,
			Age:         data.Age,
			Gender:      data.Gender,
			Nationality: data.Nationality,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
func EnrichPerson(p model.Person) model.Person {
	name := p.Name
	// Получаем возраст
	if age := getAge(name); age != nil {
		p.Age = age
	}
	// Получаем пол
	if age := getGender(name); age != nil {
		p.Gender = age
	}
	// Получаем национальность
	if nat := getNationality(name); nat != nil {
		p.Nationality = nat
	}
	return p
}

func getAge(name string) *int {
	url := fmt.Sprintf("https://api.agify.io/?name=%s", name)
	body := fetch(url)
	var res struct {
		Age int `json:"age"`
	}
	if err := json.Unmarshal(body, &res); err == nil && res.Age > 0 {
		return &res.Age
	}
	return nil
}

func getGender(name string) *string {
	url := fmt.Sprintf("https://api.genderize.io/?name=%s", name)
	body := fetch(url)
	var res struct {
		Gender string `json:"gender"`
	}
	if err := json.Unmarshal(body, &res); err == nil && res.Gender != "" {
		return &res.Gender
	}
	return nil
}

func getNationality(name string) *string {
	url := fmt.Sprintf("https://api.nationalize.io/?name=%s", name)
	body := fetch(url)
	var res struct {
		Country []struct {
			CountryID string  `json:"country_id"`
			Prob      float64 `json:"probability"`
		} `json:"country"`
	}
	if err := json.Unmarshal(body, &res); err == nil && len(res.Country) > 0 {
		return &res.Country[0].CountryID
	}
	return nil
}

// fetch выполняет простой запрос GET-запрос и возвращает тело ответа
func fetch(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return body
}
