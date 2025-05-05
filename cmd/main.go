package main

import (
	"effective-mobail/internal/config"
	"effective-mobail/internal/handler"
	"effective-mobail/internal/storage"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	cfg := config.LoadConfig() // Загружаю конфиг из .env

	// Инициализирую подключение к БД
	db := storage.InitPostgres(cfg)
	h := handler.NewPeopleHandler(db)
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Сервер запущен и работает"))
	})

	// REST API
	r.HandleFunc("/people", h.CreatePerson).Methods("POST")
	r.HandleFunc("/people", h.GetPeople).Methods("GET")
	r.HandleFunc("/people/{id}", h.UpdatePerson).Methods("PUT")
	r.HandleFunc("/people/{id}", h.DeletePerson).Methods("DELETE")

	log.Println("Сервер запущен на порту:", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal("не удалось запустить сервер: %v", err)
	}
}
