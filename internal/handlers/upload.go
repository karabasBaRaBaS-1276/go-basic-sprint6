package handlers

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

// Структура, облуживающая запросы по файлам.
type FileHandler struct {
	fileFieldName string // Имя поля, в котором поступает файл для загрузки
}

// NewFileHandler позволяет получить новый экземпляр [FileHandler].
func NewFileHandler() *FileHandler {
	return &FileHandler{fileFieldName: "myFile"}
}

// Обработка эндпоинта '/upload'.
// Допустим только POST метод для загрузки файла
func (fileHandler *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	log.Printf("Поступил запрос %s: '%s'\n", r.Method, r.URL.String())

	if r.Method != http.MethodPost {
		http.Error(w, fmt.Sprintf("%s method not allowed", r.Method), http.StatusInternalServerError) // хотя больше подходит http.StatusMethodNotAllowed
	}

	// Разбираем, что к нам пришло
	r.ParseMultipartForm(10 << 20) // ограничение не более 10 Мб

	fileInput, handler, err := r.FormFile(fileHandler.fileFieldName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// отложенное закрытие файла
	defer fileInput.Close()

	log.Printf("На обработку поступил файл: '%s'; Размер: %d байт; Тип файла: %s", handler.Filename, handler.Size, handler.Header.Get("Content-Type"))

	// Надо бы проверить, что данные в файле - это текст, который подходит для наших задач.
	// По идее, можно прочитать первые несколько байт и проверить, являются ли они ASCII-символами
	isTextData, err := isTextFile(fileInput)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isTextData {
		http.Error(w, "Передан не текстовый файл. Попробуйте другой файл", http.StatusInternalServerError)
		return
	}

	// Подготовим файл для записи
	fileResultPath := getFileNameOutput(filepath.Ext(handler.Filename))
	log.Printf("Имя файла с результатом: %s", fileResultPath)
	fileResult, err := os.Create(fileResultPath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer fileResult.Close()

	// создаем сканер для чтения файла
	scanner := bufio.NewScanner(fileInput)
	// читаем файл построчно
	lineNum := 1
	lenByte := 0
	converter := service.NewConverterMorse()
	for scanner.Scan() {
		// Так как по заданию нет ограничений на размер файла, то будем обработку делать построчно
		log.Printf("Обработка строки с содержимым: %s", scanner.Text())
		res := fmt.Sprintf("%s\n", converter.Process(scanner.Text())) // Риски. Каждая строка обрабатывается независимо...
		log.Printf("Результат обработки строки: %s", res)
		// Запишем результат в локальный файл
		n, err := fileResult.WriteString(res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// и в ответ
		_, err = io.WriteString(w, res)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		lenByte += n
		lineNum++
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
}

// Вернет true если файл текстовый.
func isTextFile(fileInput multipart.File) (bool, error) {
	buf := make([]byte, 512) // Проверяем первые 512 байт
	nRows, err := fileInput.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}
	// Сбрасываем позицию чтения в начало файла
	_, err = fileInput.Seek(0, io.SeekStart)
	if err != nil {
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

// Вернет имя файла для записи результата
func getFileNameOutput(ext string) string {
	// Подготовим файл для записи
	res := time.Now().UTC().Format("2006-01-02T15-04-05.999999999")
	res = strings.ReplaceAll(res, ".", "_")
	res += ext
	return res
}
