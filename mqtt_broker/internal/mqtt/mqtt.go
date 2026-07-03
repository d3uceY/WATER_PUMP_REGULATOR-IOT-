package mqtt

import (
	"fmt"

	"mqtt_broker/internal/telegram"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func ConnectMQTT() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("go_mqtt_broker_client")

	client := mqtt.NewClient(opts)

	token := client.Connect()
	token.Wait()

	if err := token.Error(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to MQTT fr")

	return client
}

func RunSubscriptions(client mqtt.Client, telegramClient *telegram.TelegramClient) {
	topics := []string{topicPumpOff, topicPumpOn}

	for _, topic := range topics {
		token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), string(msg.Payload()))
			if err := sendTopicMessage(msg.Topic(), telegramClient); err != nil {
				fmt.Printf("Failed to send topic message: %v\n", err)
			}
		})
		token.Wait()
		if err := token.Error(); err != nil {
			panic(err)
		}
	}
}
