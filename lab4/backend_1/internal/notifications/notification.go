package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type TelegramNotification struct {
	ChatID string `json:"chat_id"`
	Text   string `json:"text"`
}

func SendTelegramNotification(webhookURL string, chatID string, message string) error {
	notification := TelegramNotification{
		ChatID: chatID,
		Text:   message,
	}

	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("Ошибка отправки уведомления: статус %d", resp.StatusCode)
	}

	return nil
}
