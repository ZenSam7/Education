package tools

import (
	"fmt"
	"gopkg.in/gomail.v2"
	"os"
	"strings"
)

type EmailSender interface {
	SendMail(toEmail, templateFile string, templateValues map[string]string) error
}

type GmailSender struct {
	Config Config
}

// SendMail отправляет письмо с использованием указанного шаблона
// (путь надо к файлу из директории откуда эта вызавается функция)
// TemplateFile Файл в папке templates
// TemplateValues Какие плэйсхолдеры на что заменяем (можно вставлять числа)
func (m *GmailSender) SendMail(toEmail, templateFile string, templateValues map[string]string) error {
	// Новый ящик для отправки писем
	msg := gomail.NewMessage()
	msg.SetHeader("From", m.Config.EmailSender)
	msg.SetHeader("Subject", "Верифицируйте почту Education")

	dialer := gomail.NewDialer(m.Config.EmailHost, m.Config.EmailPort, m.Config.EmailSender, m.Config.EmailPassword)

	// Чтение шаблона письма из файла
	templateData, err := os.ReadFile(strings.Join([]string{absolutePath, "templates", templateFile}, PathSeparator))
	if err != nil {
		return fmt.Errorf("не удалось прочитать файл шаблона: %v", err)
	}

	// Замена placeholder в шаблоне на реальные значения
	body := string(templateData)
	for key, val := range templateValues {
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
