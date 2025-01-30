package device

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Pc struct {
	Device
	Power bool `json:"power"`
	mu    sync.RWMutex
}

func (pc *Pc) UpdatePower(power bool) {
	pc.mu.Lock()
	defer pc.mu.Unlock()
	pc.Power = power
}

func (pc *Pc) ListenUpdate() {
	client := *pc.client
	topic := pc.Device.getTopicUpdate()
	client.Subscribe(topic, 0, func(c mqtt.Client, msg mqtt.Message) {
		var data map[string]string
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Printf("Invalid command: %v", err)
			return
		}
		value, ok := data["power"]
		if !ok {
			log.Printf("Invalid data: %s", value)
			return
		}
		bool, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Invalid data: %s", value)
			return
		}

		pc.UpdatePower(bool)
		pc.ChangeState(IDLE)
	})
}

func (pc *Pc) GetDataResp() map[string]any {
	pc.mu.RLock()
	defer pc.mu.RUnlock()
	return map[string]any{
		"id":         pc.Id,
		"name":       pc.Name,
		"type":       pc.Type.String(),
		"state":      pc.Device.State.String(),
		"last_check": pc.Device.LastCheck,
		"power":      pc.Power,
	}
}

func (pc *Pc) UnListenUpdate() {
	client := *pc.client
	topic := pc.Device.getTopicUpdate()
	client.Unsubscribe(topic)
}

func (pc *Pc) RequestToggle() error {
	if !pc.Device.IsState(IDLE) {

		return fmt.Errorf("device on state %s is unavailable", pc.Device.State.String())
	}
	pc.ChangeState(PROCESSING)
	return nil
}
