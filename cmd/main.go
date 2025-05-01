package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Простая маршрутизация
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Сервер запущен и работает"))

	})

	log.Println("Сервер запущен на порту:", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("не удалось запустить сервер: %v", err)
	}
}
