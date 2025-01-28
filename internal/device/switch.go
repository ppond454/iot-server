package device

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Switch struct {
	Device
	isOn bool
}

func (sw *Switch) ListenUpdate(onUpdate func(), client mqtt.Client) {
	topic := fmt.Sprintf("sw/%s/update", sw.Id)
	client.Subscribe(topic, 0, func(c mqtt.Client, msg mqtt.Message) {
		var data map[string]string
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Printf("Invalid command: %v", err)
			return
		}
		value, ok := data["isOn"]
		if !ok {
			log.Printf("Invalid data: %s", value)
			return
		}
		bool, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Invalid data: %s", value)
			return
		}

		sw.isOn = bool
		onUpdate()
	})
}
