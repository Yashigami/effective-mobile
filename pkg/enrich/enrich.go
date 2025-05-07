package enrich

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Структура ответа от API возраста
type AgeResponse struct {
	Age int `json:"age"` // Возраст
}

// Структура ответа от API пола
type GenderResponse struct {
	Gender string `json:"gender"`
}

// Структура ответа от API национальности
type NationalityResponse struct {
	Country []struct {
		Country     string  `json:"country_id"`  // Код страны
		Probability float64 `json:"probability"` // Вероятность
	} `json:"country"`
}

// Структура, объединяющая все вместе
type EnrichedData struct {
	Age         int
	Gender      string
	Nationality string
}

// Функция обогащения
func Enrich(name string) (*EnrichedData, error) {
	client := &http.Client{Timeout: 10 * time.Second} // Ограничение по времени на каждый запрос

	var (
		ageResp         AgeResponse
		genderResp      GenderResponse
		nationalityResp NationalityResponse
	)

	// 1. Получаем возраст
	ageURL := fmt.Sprintf("https://api.agify.io?name=%s", name)
	if err := fetchJSON(client, ageURL, &ageResp); err != nil {
		return nil, fmt.Errorf("ошибка при получении возраста: %w", err)
	}

	// 2. Получаем пол
	genderURL := fmt.Sprintf("https://api.genderize.io?name=%s", name)
	if err := fetchJSON(client, genderURL, &genderResp); err != nil {
		return nil, fmt.Errorf("ошибка при получении пола: %w", err)
	}

	// 3. Получаем национальность
	natUrl := fmt.Sprintf("https://api.nationalize.io?name=%s", name)
	if err := fetchJSON(client, natUrl, &nationalityResp); err != nil {
		return nil, fmt.Errorf("ошибка при получении национальности: %w", err)
	}

	// Берем первую страну
	nationality := ""
	if len(nationalityResp.Country) > 0 {
		nationality = nationalityResp.Country[0].Country
	}

	return &EnrichedData{
		Age:         ageResp.Age,
		Gender:      genderResp.Gender,
		Nationality: nationality,
	}, nil
}

// Вспомогательная функция для запроса и декодинга JSON
func fetchJSON(client *http.Client, url string, target interface{}) error {
	resp, err := client.Get(url)
	if err != nil {
		return err // Ошибка при подключении
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("неуспешный ответ: %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(target) // Декодим JSON в структуру
}
