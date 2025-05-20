package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Структура, облуживающая запросы на получения статичных данных.
type StaticHandler struct {
	htmlIndexFile string
}

// NewStaticHandler позволяет получить новый экземпляр [StaticHandler].
func NewStaticHandler() *StaticHandler {
	return &StaticHandler{htmlIndexFile: "index.html"}
}

// Обработка энпоинта '/'
// Допустим только GET метод
func (staticHandler *StaticHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	defer func() {
		log.Printf("Обработан запрос %s: '%s' (время обработки: %s)\n", r.Method, r.URL.String(), time.Since(start))
	}()

	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("%s method not allowed", r.Method), http.StatusInternalServerError) // хотя больше подходит http.StatusMethodNotAllowed
	}

	file, err := os.Open(staticHandler.htmlIndexFile)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	// копируем содержимое файла в ответ
	_, err = io.Copy(w, file)
	if err != nil {
		http.Error(w, "ошибка при отправке файла", http.StatusInternalServerError)
		return
	}

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
