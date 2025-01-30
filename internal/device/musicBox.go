package device

import (
	"encoding/json"
	"log"
	"strconv"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MusicBox struct {
	Device
	IsOn      bool
	IsPlaying bool
}

func (m *MusicBox) Get() *MusicBox {
	m.Device.mu.Lock()
	defer m.Device.mu.Unlock()

	return m
}

func (m *MusicBox) SetPlaying(isPlaying bool) *MusicBox {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.IsPlaying = isPlaying
	return m
}

func (m *MusicBox) Play() *MusicBox {
	// m.isPlaying = true
	return m
}

func (m *MusicBox) ListenUpdate() {
	topic := m.Device.getTopicUpdate()
	client := *m.client
	client.Subscribe(topic, 0, func(c mqtt.Client, msg mqtt.Message) {
		var data map[string]string
		if err := json.Unmarshal(msg.Payload(), &data); err != nil {
			log.Printf("Invalid command: %v", err)
			return
		}
		value, ok := data["isPlaying"]
		if !ok {
			log.Printf("Invalid data: %s", value)
			return
		}
		bool, err := strconv.ParseBool(value)
		if err != nil {
			log.Printf("Invalid data: %s", value)
			return
		}

		m.SetPlaying(bool)
	})
}

func (m *MusicBox) RequestTurnOn() {
}

func (m *MusicBox) RequestTurnOff() {

}

func (m *MusicBox) RequestToPlay() {

}
