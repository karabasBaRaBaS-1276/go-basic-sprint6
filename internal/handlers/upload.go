package handlers

import (
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
)

// Структура, облуживающая запросы по файлам.
type FileHandler struct {
	fileFieldName string // Имя поля, в котором поступает файл для загрузки
}

// NewFileHandler позволяет получить новый экземпляр [FileHandler].
func NewFileHandler() *FileHandler {
	return &FileHandler{fileFieldName: "myFile"}
}

// Обработка энпоинта '/upload'.
// Допустим только POST метод для загрузки файла
func (fileHandler *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("Поступил запрос %s: '%s'\n", r.Method, r.URL.String())

	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("%s method not allowed", r.Method), http.StatusInternalServerError) // хотя больше подходит http.StatusMethodNotAllowed
	}

	// Разбираем, что к нам пришло
	r.ParseMultipartForm(10 << 20) // ограничение не более 10 Мб

	file, handler, err := r.FormFile(fileHandler.fileFieldName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// отложенное закрытие файла
	defer file.Close()

	log.Printf("На обработку поступил файл: '%s'; Размер: %d байт; Тип файла: %s", handler.Filename, handler.Size, handler.Header.Get("Content-Type"))

	// Надо бы проверить, что данные в файле - это текст, который подходит для наших задач.
	// По идее, можно прочитать первые несколько байт и проверить, являются ли они ASCII-символами
	isTextData, err := isTextFile(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isTextData {
		http.Error(w, "Передан не текстовый файл. Попробуйте другой файл", http.StatusInternalServerError)
		return
	}

	// Записываем содержимое файла в массив строк

	w.Header().Set("Content-Type", "application/json")
	http.Error(w, `{"error": "Обработка работы с файлом пока не реализована"}`, http.StatusMethodNotAllowed)
}

// Вернет true если файл текстовый.
func isTextFile(file multipart.File) (bool, error) {
	buf := make([]byte, 512) // Проверяем первые 512 байт
	nRows, err := file.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}
	log.Printf("Первые байты файла (не более 512):\n%s", buf)

	// Проверяем, что нет нулевых байт и большинство символов - печатные
	for i := 0; i < nRows; i++ {
		if buf[i] == 0 {
			log.Printf("Обнаружен бинарный файл: 0x00 в позиции %d", i)
			return false, nil // Бинарные файлы часто содержат 0x00
		}
		if buf[i] < 32 && buf[i] != '\n' && buf[i] != '\r' && buf[i] != '\t' {
			log.Printf("Найден не текстовый байт в позиции %d: 0x%02x (dec %d --> %s)", i, buf[i], buf[i], string(buf[i]))
			return false, nil // Непечатные символы (кроме \n, \r, \t)
		}
		// Всё остальное (включая UTF-8) считаем текстом
	}
	log.Printf("Первые %d байт файла выглядят как текстовые :)", nRows)
	return true, nil
}
