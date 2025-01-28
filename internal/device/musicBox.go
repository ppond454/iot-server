package device

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"sync"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MusicBox struct {
	Device
	isPlaying bool
	mu        sync.Mutex
}

func (m *MusicBox) Get() *MusicBox {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m
}

func (m *MusicBox) ChangeState(isPlaying bool) *MusicBox {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.isPlaying = isPlaying
	return m
}

func (m *MusicBox) Play() *MusicBox {
	// m.isPlaying = true
	return m
}

func (m *MusicBox) ListenUpdate(onUpdate func(), client mqtt.Client) {
	topic := fmt.Sprintf("music_box/%s/update", m.Id)
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

		m.ChangeState(bool)
		onUpdate()
	})
}

func (m *MusicBox) UnListenUpdate(client mqtt.Client) {
	topic := fmt.Sprintf("music_box/%s/update", m.Id)
	client.Unsubscribe(topic)
}
