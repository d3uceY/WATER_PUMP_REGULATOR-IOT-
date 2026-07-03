package mqtt

import (
	"errors"

	"mqtt_broker/internal/telegram"
	"mqtt_broker/internal/whatsapp"
)

const (
	topicPumpOn  = "message_pump_on"
	topicPumpOff = "message_pump_off"
)

func sendTopicMessage(topic string, telegramClient *telegram.TelegramClient) error {
	var message string

	switch topic {
	case topicPumpOn:
		message = "pump is turned on"
	case topicPumpOff:
		message = "pump is turned off"
	default:
		return nil
	}

	return errors.Join(
		whatsapp.SendMessage(message),
		telegramClient.SendToAllChats(message),
	)
}
