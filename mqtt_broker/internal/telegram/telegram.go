package telegram

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"mqtt_broker/internal/config"
	store "mqtt_broker/internal/storage"
)

type GetUpdatesResponse struct {
	OK     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type Update struct {
	UpdateID int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	MessageID int    `json:"message_id"`
	Date      int64  `json:"date"`
	Text      string `json:"text"`
	Chat      Chat   `json:"chat"`
	From      User   `json:"from"`
}

type Chat struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type User struct {
	ID           int64  `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type TelegramClient struct {
	Token      string
	UpdatesURI string
	PostURI    string
}

// i think this one is obvious, bro
func (t *TelegramClient) InitTelegram() {
	t.Token = config.Get("TELEGRAM_BOT_TOKEN")
	t.PostURI = fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", t.Token)
	t.UpdatesURI = fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", t.Token)
}

// this sends the message to the users subscribed to the bot
// by taking the message and chatid of the users
func (t *TelegramClient) Send(message string, chatId int64) {

	payload := map[string]any{
		"chat_id": chatId,
		"text":    message,
	}

	body, _ := json.Marshal(payload)

	resp, err := http.Post(t.PostURI, "application/json", bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println(resp.Status)
}

// this gets the chat ids
func (t *TelegramClient) GetChatIDs() ([]int64, error) {
	updates, err := GetUpdates(t.Token)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	for _, update := range updates.Result {
		chatID := update.Message.Chat.ID
		if chatID == 0 || store.ChatIdExists(int64(chatID)) {
			continue
		}

		// i will insert the necessary data into the db here
		data := store.TelegramChat{
			ChatId:   int64(chatID),
			Username: update.Message.Chat.Username,
		}
		store.InsertChat(data)
	}

	chats, err := store.GetChats()

	if err != nil {
		return nil, err
	}

	chatIDs := make([]int64, len(chats))

	for _, chat := range chats {
		chatIDs = append(chatIDs, chat.ChatId)
	}

	return chatIDs, nil
}

// this sends to all chatIds using a go routine
// 'cause i am all about that concurrency and shit
func (t *TelegramClient) SendToAllChats(message string) error {
	chatIDs, err := t.GetChatIDs()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, chatID := range chatIDs {
		wg.Add(1)
		go func(chatID int64) {
			defer wg.Done()
			t.Send(message, chatID)
		}(chatID)
	}
	wg.Wait()

	return nil
}

// i use this to get
// the updates from the bot and
// get the chat ids of the users
func GetUpdates(token string) (*GetUpdatesResponse, error) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/getUpdates", token)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var updates GetUpdatesResponse
	if err := json.NewDecoder(resp.Body).Decode(&updates); err != nil {
		return nil, err
	}

	return &updates, nil
}
