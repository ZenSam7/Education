package validator

import (
	"fmt"
	"net/mail"
)

type anyNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64 |
		~complex64 | ~complex128
}

// ValidateString Контролируем размер строки
// minLenght по умолчанию 0, maxLenght по умолчанию ∞
func ValidateString(s string, minLenght, maxLenght int) error {
	if len(s) < minLenght || (maxLenght != 0 && len(s) > maxLenght) {
		return fmt.Errorf("строка должна содержать %d-%d символов", minLenght, maxLenght)
	}

	return nil
}

func ValidateNaturalNum(x int) error {
	if x <= 0 {
		return fmt.Errorf("число должно быть натуральным")
	}

	return nil
}

func ValidateEmail(s string) error {
	if err := ValidateString(s, 5, 0); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(s); err != nil {
		return fmt.Errorf("неправильный адрес электронной почты")
	}

	return nil
}
