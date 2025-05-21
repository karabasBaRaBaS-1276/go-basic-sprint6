package service

import (
	"unicode"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/pkg/morse"
)

// Структура конвертер азбуки Морзе
type ConverterMorse struct {
}

func NewConverterMorse() *ConverterMorse {
	return &ConverterMorse{}
}

// Обработка строки.
//
// Если передан обычный текст, метод должен переконвертировать его в код Морзе и вернуть;
// и наоборот — если был передан код Морзе, метод должен переконвертировать его в обычный текст и вернуть.
func (converter *ConverterMorse) Process(line string) string {

	var result string

	if isMorseString(line) {
		result = morse.ToText(line)
	} else {
		result = morse.ToMorse(line)
	}
	return result
}

// Вернет true, если переданная строка состоит из набора символов азбуки Морзе
func isMorseString(s string) bool {
	for _, r := range s {
		if r != '.' && r != '-' && !unicode.IsSpace(r) {
			return false
		}
	}
	return true
}
