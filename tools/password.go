package tools

import "golang.org/x/crypto/bcrypt"

// GetPasswordHash Генерируем хеш
func GetPasswordHash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// CheckPassword Проверяем пароль
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GetRandomHash Генерируем случайный хеш
func GetRandomHash() string {
	hash, _ := GetPasswordHash(GetRandomString())
	return hash
}
