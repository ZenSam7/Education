package tools

import (
	"fmt"
	"net/mail"
)

// ValidateString Контролируем размер строки
// minLenght по умолчанию 0, maxLenght по умолчанию ∞
func ValidateString(s string, minLenght, maxLenght int) error {
	if len(s) < minLenght || (maxLenght != 0 && len(s) > maxLenght) {
		if maxLenght != 0 {
			return fmt.Errorf("строка должна содержать от %d до %d символов", minLenght, maxLenght)
		}
		return fmt.Errorf("строка должна содержать минимум %d символов", minLenght)
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
