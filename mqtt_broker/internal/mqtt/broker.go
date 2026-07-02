package mqtt

import (
	"fmt"

	mochi "github.com/mochi-mqtt/server/v2"
	"github.com/mochi-mqtt/server/v2/hooks/auth"
	"github.com/mochi-mqtt/server/v2/listeners"
)

func StartBroker() {
	server := mochi.New(nil)

	_ = server.AddHook(new(auth.AllowHook), nil)

	tcp := listeners.NewTCP(listeners.Config{
		ID:      "tcp",
		Address: ":1883",
	})
	if err := server.AddListener(tcp); err != nil {
		panic(err)
	}

	go func() {
		if err := server.Serve(); err != nil {
			panic(err)
		}
	}()

	fmt.Println("MQTT broker listening on :1883")
}
