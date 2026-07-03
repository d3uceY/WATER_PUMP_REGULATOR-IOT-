package main

import (
	"fmt"
	"net/http"

	"mqtt_broker/internal/config"
	"mqtt_broker/internal/mqtt"
	"mqtt_broker/internal/telegram"
	"mqtt_broker/internal/whatsapp"
)

func main() {
	config.Load()

	if err := whatsapp.Init(); err != nil {
		panic(err)
	}

	telegramClient := telegram.TelegramClient{}
	telegramClient.InitTelegram()

	mqtt.StartBroker()

	client := mqtt.ConnectMQTT()
	defer client.Disconnect(250)

	mqtt.RunSubscriptions(client)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running. no cap on god, bro")
	})

	fmt.Println("Listening on :8080")

	http.ListenAndServe(":8080", nil)
}
