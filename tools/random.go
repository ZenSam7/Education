package tools

import (
	"math/rand"
)

// GetRandomString minLength: необязательный аргумент. По умолчанию минимальная длина: 2
func GetRandomString(minLength ...int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Установка минимальной длины строки
	var minLengthString int
	if len(minLength) > 0 {
		minLengthString = minLength[0]
	} else {
		minLengthString = 2
	}

	// Генерация случайной длины строки от minLengthString до minLengthString + len(letters)
	strLength := rand.Intn(len(letters)) + minLengthString

	// Создание пустой строки заданной длины
	str := make([]byte, strLength)
	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}

// GetRandomInt Число может быть отрицательным
func GetRandomInt() int32 {
	return rand.Int31() * int32(1-2*rand.Intn(2))
}

// GetRandomUint Число не может быть отрицательным
func GetRandomUint() int32 {
	return rand.Int31() + 1
}

// GetRandomEmail Генерируем случайную почту
func GetRandomEmail() string {
	return GetRandomString() + "@" + GetRandomString() + ".com"
}
