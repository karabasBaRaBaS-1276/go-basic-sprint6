package handlers

import "net/http"

// Структура, облуживающая запросы на получения статичных данных.
type StaticHandler struct{}

// NewStaticHandler позволяет получить новый экземпляр [StaticHandler].
func NewStaticHandler() *StaticHandler {
	return &StaticHandler{}
}

func (*StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"error": "Обработка корневого эндпоинта пока не реализована"}`, http.StatusMethodNotAllowed)
}

// Структура, облуживающая запросы по файлам.
type FileHandler struct{}

// NewFileHandler позволяет получить новый экземпляр [FileHandler].
func NewFileHandler() *FileHandler {
	return &FileHandler{}
}

func (*FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"error": "Обработка работы с файлом пока не реализована"}`, http.StatusMethodNotAllowed)
}
