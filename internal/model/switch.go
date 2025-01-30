package model

import (
	"encoding/json"
	"log"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type ToggleSwitch struct {
	Device
	IsOn bool
}

func (sw *ToggleSwitch) ListenUpdate() {
	topic := sw.Device.getTopicUpdate()
	client := *sw.client
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

		sw.IsOn = bool
	})
}

func (sw *ToggleSwitch) RequestTurnOn() {

}

func (sw *ToggleSwitch) RequestTurnOff() {

}
