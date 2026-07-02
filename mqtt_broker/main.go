package main

import (
	"fmt"
	"net/http"
	"mqtt_broker/internal/config"
	"mqtt_broker/internal/mqtt"
)

func main() {
	config.Load()

	client := mqtt.ConnectMQTT()
	defer client.Disconnect(250)

	mqtt.RunSubscriptions(client)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running. no cap on god, bro")
	})

	fmt.Println("Listening on :8080")

	http.ListenAndServe(":8080", nil)
}