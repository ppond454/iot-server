package device

import (
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type IoTDevice interface {
	IsState(state State) bool
	GetState() State
	ChangeState(state State)
	GetDataResp() map[string]any
	GetData() Device
	ListenUpdate()
	UnListenUpdate()
	RequestToggle() error
	SetLastCheck(time *time.Time)
}

type Device struct {
	Id        string     `json:"id"`
	Name      string     `json:"name"`
	Type      DeviceType `json:"type"` // PC, SWITCH, MUSIC_BOX, etc.
	State     State      `json:"state"`
	LastCheck *time.Time `json:"last_check,omitempty"`
	client    *mqtt.Client
	callback  func(*Device)
	mu        sync.RWMutex
}

func NewDevice(_type, id, name string, client *mqtt.Client, callback func(*Device)) IoTDevice {
	switch _type {
	case "PC":
		return &Pc{
			Device: Device{
				Id:        id,
				Name:      name,
				Type:      PC,
				State:     DISCONNECTED,
				LastCheck: nil,
				client:    client,
				callback:  callback,
			},
			Power: false,
		}
	// case "TOGGLE_SWITCH":
	// 	return &ToggleSwitch{Device: Device{id, name, TOGGLE_SWITCH, DISCONNECTED, nil, client}, IsOn: false}
	// case "MUSIC_BOX":
	// 	return &MusicBox{Device: Device{id, name, MUSIC_BOX, DISCONNECTED, nil, client}, IsPlaying: false, IsOn: false}
	default:
		return nil
	}
}

func (d *Device) ChangeState(state State) {
	d.mu.Lock()
	d.State.changeState(state)
	d.mu.Unlock()
	d.callback(d)
}

func (d *Device) SetLastCheck(time *time.Time) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.LastCheck = time
}

func (d *Device) IsState(state State) bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.State == state
}

func (d *Device) GetState() State {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.State
}

func (d *Device) getTopicUpdate() string {
	return fmt.Sprintf("update/%s", d.Id)
}

func (d *Device) GetData() Device {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return Device{
		Id:        d.Id,
		Name:      d.Name,
		Type:      d.Type,
		State:     d.State,
		LastCheck: d.LastCheck,
	}
}

func (d *Device) UnListenUpdate() {
	client := *d.client
	topic := d.getTopicUpdate()
	client.Unsubscribe(topic)
}
