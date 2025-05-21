package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

// Структура, облуживающая запросы на получения статичных данных.
type StaticHandler struct {
	htmlIndexFile string // Имя index файла для отдачи клиенту
}

// NewStaticHandler позволяет получить новый экземпляр [StaticHandler].
func NewStaticHandler() *StaticHandler {
	return &StaticHandler{htmlIndexFile: "index.html"}
}

// Обработка эндпоинта '/'
// Допустим только GET метод
func (staticHandler *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("Поступил запрос %s: '%s'\n", r.Method, r.URL.String())

	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%s method not allowed", r.Method), http.StatusInternalServerError) // хотя больше подходит http.StatusMethodNotAllowed
	}

	file, err := os.Open(staticHandler.htmlIndexFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// отложенное закрытие файла
	defer file.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// копируем содержимое файла в ответ
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "ошибка при отправке файла", http.StatusInternalServerError)
		return
	}

}
