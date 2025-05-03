package main

import (
	"effective-mobail/internal/config"
	"effective-mobail/internal/handler"
	"effective-mobail/internal/storage"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig() // Загружаю конфиг из .env

	// Инициализирую подключение к БД
	db := storage.InitPostgres(cfg)
	h := handler.NewPeopleHandler(db)

	// Простая маршрутизация
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Сервер запущен и работает"))

	})

	http.HandleFunc("/people", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			h.CreatePerson(w, r)
		case http.MethodGet:
			h.GetPeople(w, r)
		default:
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Сервер запущен на порту:", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, nil)
	if err != nil {
		log.Fatal("не удалось запустить сервер: %v", err)
	}
}
