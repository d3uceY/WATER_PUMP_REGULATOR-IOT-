package whatsapp

import (
	"context"
	"fmt"

	"mqtt_broker/internal/config"

	whatsapp "github.com/KARTIKrocks/gowhatsapp"
)

func SendMessage(to string, message string) error {
	client, err := whatsapp.New(whatsapp.Config{
		PhoneNumberID: config.Get("WHATSAPP_PHONE_NUMBER_ID"),
		AccessToken:   config.Get("WHATSAPP_ACCESS_TOKEN"),
	})
	if err != nil {
		return fmt.Errorf("failed to create whatsapp client: %w", err)
	}

	ctx := context.Background()

	res, err := client.SendText(ctx, to, message)
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	fmt.Printf("Sent message %s to %s\n", res.MessageID, to)
	return nil
}
