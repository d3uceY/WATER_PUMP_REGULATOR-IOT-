package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func ConnectMQTT() mqtt.Client {
	opts := mqtt.NewClientOptions()
	opts.AddBroker("tcp://localhost:1883")
	opts.SetClientID("my_mqtt_client")

	client := mqtt.NewClient(opts)

	token := client.Connect()
	token.Wait()

	if err := token.Error(); err != nil {
		panic(err)
	}

	fmt.Println("Connected to MQTT fr")

	return client
}

func RunSubscriptions(client mqtt.Client) {
	topics := []string{"message_tank_off", "message_tank_on"}

	for _, topic := range topics {
		token := client.Subscribe(topic, 0, func(client mqtt.Client, msg mqtt.Message) {
			fmt.Printf("Received message on topic %s: %s\n", msg.Topic(), string(msg.Payload()))
		})
		token.Wait()
		if err := token.Error(); err != nil {
			panic(err)
		}
	}
}