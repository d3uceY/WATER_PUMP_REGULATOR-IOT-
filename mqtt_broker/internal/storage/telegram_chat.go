package storage

import (
	"fmt"
)

type TelegramChat struct {
	Id       *int
	ChatId   int64
	Username string
}

func GetChats() ([]TelegramChat, error) {
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

func InsertChat(data TelegramChat) error {
	query := `INSERT INTO telegram_chats (username, chat_id) VALUES (?, ?)`

	_, err := DB.Exec(query, data.Username, data.ChatId)

	if err != nil {
		return fmt.Errorf("failed to insert Chat: %v", err)
	}
	return nil
}

func ChatIdExists(chatId int64) bool {
	query := `SELECT 1 FROM telegram_chats WHERE chat_id = ? LIMIT 1`

	var exists int
	err := DB.QueryRow(query, chatId).Scan(&exists)

	return err == nil
}