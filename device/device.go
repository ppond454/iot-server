package device

import (
	"time"
)

type Device struct {
	Id          string
	Name        string
	isConnected bool
	lastCheck   *time.Time
}

func NewDevice(id, name string) *Device {
	return &Device{id, name, false, nil}
}

func (d *Device) IsConnected() bool {
	return d.isConnected
}

func (d *Device) Connected(time *time.Time) {
	d.isConnected = true
	d.lastCheck = time
}

func (d *Device) Disconnect() {
	d.isConnected = false
}

func (d *Device) Ping() (float32, error) {
	// Implement device health check
	return 0, nil
}

func (d *Device) GetData() Device {
	// fmt.Printf("\nUpdated Device: %+v\n", d)
	// fmt.Printf("Last Check: %s\n", d.lastCheck.Format(time.RFC1123))
	return *d
}
