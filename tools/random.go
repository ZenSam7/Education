package tools

import "math/rand"

func GetRandomString() string {
	const letters = " abcdefghijklmnopqrstuvwxyz ABCDEFGHIJKLMNOPQRSTUVWXYZ"

	// Минимальная длина: 2
	str := make([]byte, rand.Intn(len(letters)-2)+2)
	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}

// GetRandomInt Число может быть отрицательным
func GetRandomInt() int32 {
	return rand.Int31() * int32(1-2*rand.Intn(2))
}
