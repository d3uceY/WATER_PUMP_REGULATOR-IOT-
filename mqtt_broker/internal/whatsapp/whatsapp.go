package whatsapp

import (
	"context"
	"fmt"

	"mqtt_broker/internal/config"

	whatsapp "github.com/KARTIKrocks/gowhatsapp"
)

var (
	to     = config.Get("WHATSAPP_RECIPIENT_NUMBER")
	client *whatsapp.Client
)

type TopicsType struct {
	PumpOn  string
	PumpOff string
}

var topics = TopicsType{
	PumpOn:  "message_pump_on",
	PumpOff: "message_pump_off",
}

func Init() error {
	var err error
	client, err = whatsapp.New(whatsapp.Config{
		PhoneNumberID: config.Get("WHATSAPP_PHONE_NUMBER_ID"),
		AccessToken:   config.Get("WHATSAPP_ACCESS_TOKEN"),
	})
	if err != nil {
		return fmt.Errorf("failed to create whatsapp client: %w", err)
	}
	return nil
}

func sendTopicMessage(topic string) {
	switch topic {

	case topics.PumpOn:
		message := "pump is turned on"
		SendMessage(message)

	case topics.PumpOff:
		message := "pump is turned off"
		SendMessage(message)

	}
}

func SendMessage(message string) error {
	ctx := context.Background()

	res, err := client.SendText(ctx, to, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Printf("Sent message %s to %s\n", res.MessageID, to)
	return nil
}
