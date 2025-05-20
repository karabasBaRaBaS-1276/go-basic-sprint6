package main

import (
	"log"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/server"
)

func main() {
	log := log.Default()
	server := server.Get(log)

	log.Println("Запускаем сервер")
	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
