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

	sender := GmailSender{Config: config}
	err := sender.SendMail(config.EmailSender, "email_verify.html", map[string]string{})
	require.NoError(t, err)
}
