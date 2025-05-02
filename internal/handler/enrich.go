package handler

import (
	"effective-mobail/internal/model"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

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
