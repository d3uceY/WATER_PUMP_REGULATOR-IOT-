package storage

import (
	"fmt"
)

type TelegramChat struct {
	Id       *int
	ChatId   int
	Username string
}

func GetChatIds() ([]TelegramChat, error) {
	query := `SELECT id, chat_id, username FROM telegram_chats`
	rows, err := DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var telegramChatIds []TelegramChat

	for rows.Next() {
		chat := TelegramChat{}
		err := rows.Scan(&chat.Id, &chat.ChatId, &chat.Username)

		if err != nil {
			return nil, err
		}
		telegramChatIds = append(telegramChatIds, chat)
	}
	return telegramChatIds, nil
}

func InsertChatIds(data TelegramChat) error {
	query := `INSERT INTO telegram_chats (username, chat_id) VALUES (?, ?)`

	_, err := DB.Exec(query, data.Username, data.ChatId)

	if err != nil {
		return fmt.Errorf("failed to insert clip: %v", err)
	}
}
