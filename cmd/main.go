package main

import (
	_ "effective-mobail/docs"
	"effective-mobail/internal/config"
	"effective-mobail/internal/handler"
	"effective-mobail/internal/storage"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	"log"
	"net/http"
)

// @title People API
// @version 1.0
// @description API-сервис для обогащения управления данными пользователя.
// @host localhost:8080
// @BasePath /

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
	r.Use(loggingMiddleware)

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Сервер запущен на порту:", cfg.Port)
	err := http.ListenAndServe(":"+cfg.Port, r)
	if err != nil {
		log.Fatal("не удалось запустить сервер: %v", err)
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.RequestURI, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}
