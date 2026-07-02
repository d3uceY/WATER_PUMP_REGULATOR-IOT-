package main

import (
	"fmt"
	"net/http"
	"mqtt_broker/internal/mqtt"
)

func main() {
	client := mqtt.ConnectMQTT()
	defer client.Disconnect(250)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Server is running. no cap on god, bro")
	})

	fmt.Println("Listening on :8080")

	http.ListenAndServe(":8080", nil)
}