package model

import (
	"encoding/json"
	"fmt"
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
		var body PcUpdateBody
		if err := json.Unmarshal(msg.Payload(), &body); err != nil {
			fmt.Println("Invalid JSON")
			return
		}
		pc.UpdatePower(body.Power)
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

type PcUpdateBody struct {
	Power bool `json:"power"`
}
