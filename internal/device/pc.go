package device

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Pc struct {
	Device
	isPush bool
}

func (pc *Pc) ListenUpdate(onUpdate func(), client mqtt.Client) {
	topic := fmt.Sprintf("pc/%s/update", pc.Id)
	client.Subscribe(topic, 0, func(c mqtt.Client, msg mqtt.Message) {
		var data map[string]string
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Printf("Invalid command: %v", err)
			return
		}
		value, ok := data["isPush"]
		if !ok {
			log.Printf("Invalid data: %s", value)
			return
		}
		bool, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Invalid data: %s", value)
			return
		}

		pc.isPush = bool
		onUpdate()
	})
}
