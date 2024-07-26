package tools

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"strings"
)

type EmailSender interface {
	SendMail(toEmail string) error
}

type GmailSender struct {
	Config Config
	// TemplateFile Файл в папке templates
	TemplateFile string
	// TemplateValues Какие плэйсхолдеры на что заменяем (можно вставлять числа)
	TemplateValues map[string]string
}

// SendMail отправляет письмо с использованием указанного шаблона
// (путь надо к файлу из директории откуда эта вызавается функция)
func (m *GmailSender) SendMail(toEmail string) error {
	// Новый ящик для отправки писем
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Config.EmailSender)
	msg.SetHeader("Subject", "Верифицируйте почту Education")

	dialer := gomail.NewDialer(m.Config.EmailHost, m.Config.EmailPort, m.Config.EmailUser, m.Config.EmailPassword)

	// Чтение шаблона письма из файла
	templateData, err := os.ReadFile(strings.Join([]string{absolutePath, "templates", m.TemplateFile}, PathSeparator))
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл шаблона: %v", err)
	}

	// Замена placeholder в шаблоне на реальные значения
	body := string(templateData)
	for key, val := range m.TemplateValues {
		body = strings.Replace(body, key, val, -1)
	}

	// Установка получателя и тела письма
	msg.SetHeader("To", toEmail)
	msg.SetBody("text/html", body)

	// Отправка письма
	if err := dialer.DialAndSend(msg); err != nil {
		return fmt.Errorf("не удалось отправить письмо: %v", err)
	}

	return nil
}
