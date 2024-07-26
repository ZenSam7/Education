package tools

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSendMailGmailSender(t *testing.T) {
	// Пропускаем тесты, которые не выполнятся в Github Actions
	if testing.Short() {
		t.Skip()
	}

	config := LoadConfig()

	sender := GmailSender{Config: config, TemplateFile: "email_verify.html", TemplateValues: map[string]string{}}
	err := sender.SendMail(config.EmailSender)
	require.NoError(t, err)
}
