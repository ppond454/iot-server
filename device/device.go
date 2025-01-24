package device

type Controller interface {
	Ping() (float32, error)
	Connect() (bool, error)
	IsConnected() bool
	GetData() Device
}

type Device struct {
	Id          string
	Name        string
	isConnected bool
}

func New(id, name string) Controller {
	return &Device{id, name, false}
}

func (d *Device) Connect() (bool, error) {
	d.isConnected = true
	return true, nil
}

func (d *Device) IsConnected() bool {
	return d.isConnected
}

func (d *Device) Ping() (float32, error) {
	// Implement device health check
	return 0, nil
}

func (d *Device) GetData() Device {
	return *d
}
