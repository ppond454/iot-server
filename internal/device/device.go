package device

import (
	"time"
)

type IoTDevice interface {
	IsConnected() bool
	Connected(time *time.Time)
	Disconnect()
	GetData() Device
}

type Device struct {
	Id          string
	Name        string
	isConnected bool
	LastCheck   *time.Time
}

func NewDevice(_type, id, name string) IoTDevice {
	switch _type {
	case "PC":
		return &Pc{Device: Device{id, name, false, nil}, isPush: false}
	case "SWITCH":
		return &Switch{Device: Device{id, name, false, nil}, isOn: false}
	case "MUSIC_BOX":
		return &MusicBox{Device: Device{id, name, false, nil}, isPlaying: false}
	default:
		return nil
	}
}

func (d *Device) IsConnected() bool {
	return d.isConnected
}

func (d *Device) Connected(time *time.Time) {
	d.isConnected = true
	d.LastCheck = time
}

func (d *Device) Disconnect() {
	d.isConnected = false
}

func (d *Device) GetData() Device {
	// fmt.Printf("\nUpdated Device: %+v\n", d)
	// fmt.Printf("Last Check: %s\n", d.lastCheck.Format(time.RFC1123))
	return *d
}
