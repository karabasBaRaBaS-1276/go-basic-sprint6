package server

import (
	"log"
	"net/http"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/handlers"
)

// Запускает сервер
// Принимает на вход указатель на логгер
// Возвращает:
//   - указатель на настроенный сервер для старта (*http.Server)
func Get(log *log.Logger) *http.Server {

	// Обработчики, в которых реализованы методы обработки запросов
	staticHandler := handlers.NewStaticHandler()
	fileHandler := handlers.NewFileHandler()

	// http-роутер
	mux := http.NewServeMux()
	mux.HandleFunc("/", staticHandler.ServeHTTP)     // Корневой эндпоинт
	mux.HandleFunc("/upload", fileHandler.ServeHTTP) // Эндпоинт загрузки файла

	server := &http.Server{
		Addr:         ":8080",          // порт 8080
		Handler:      mux,              // http-роутер
		ErrorLog:     log,              // указатель на логер
		ReadTimeout:  5 * time.Second,  // таймаут для чтения
		WriteTimeout: 10 * time.Second, // таймаут для записи
		IdleTimeout:  15 * time.Second, // таймаут ожидания следующего запроса
	}

	return server
}
