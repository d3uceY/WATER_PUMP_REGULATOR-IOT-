package main

import (
	"fmt"
	"net/http"

	"mqtt_broker/internal/config"
	"mqtt_broker/internal/mqtt"
	"mqtt_broker/internal/storage"
	"mqtt_broker/internal/telegram"
	"mqtt_broker/internal/whatsapp"
)

func main() {

	// sqlite db init
	storage.InitDB()

	// loads env into memory
	config.Load()

	// some whatsapp stuff
	if err := whatsapp.Init(); err != nil {
		panic(err)
	}

	// some telegram shit
	telegramClient := telegram.TelegramClient{}
	telegramClient.InitTelegram()

	// mqtt broker
	mqtt.StartBroker()

	// the funny thing is, this app is both a broker and an MQTT client at the same time
	// so basically, it is it's own client
	client := mqtt.ConnectMQTT()
	defer client.Disconnect(250)

	mqtt.RunSubscriptions(client, &telegramClient)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running. no cap on god, bro")
	})

	fmt.Println("Listening on :8080")

	http.ListenAndServe(":8080", nil)
}
