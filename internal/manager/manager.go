package manager

import (
	"errors"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	devices "github.com/ppond454/iot-backend/internal/device"
)

type IController interface {
	New(params *Params) (*List, error)
	StartAliveWorker() func()
	FindDevice(device string)
	GetDevicesJSON() map[string]any
}

type Params struct {
	Client        mqtt.Client
	OnStateChange func(*devices.Device, *List)
}

type List struct {
	devices       map[string]devices.IoTDevice
	onStateChange func(*devices.Device, *List)
	mu            sync.RWMutex
}

var client mqtt.Client = nil

func New(params *Params) (*List, error) {
	if client != nil {
		return nil, errors.New("controller already exists")
	}
	client = params.Client
	return &List{
		devices:       make(map[string]devices.IoTDevice),
		onStateChange: params.OnStateChange,
	}, nil
}

func (list *List) AddDevice(id string, d devices.IoTDevice) (map[string]devices.IoTDevice, error) {
	list.mu.Lock()
	defer list.mu.Unlock()
	if _, exist := list.devices[id]; exist {
		fmt.Printf("Device '%s' already exists\n", id)
		return nil, fmt.Errorf("device '%s' already exists", id)
	}
	list.devices[id] = d
	fmt.Printf("add device: '%s' \n", id)
	return list.devices, nil
}

func (list *List) RemoveDevice(id string) error {
	list.mu.Lock()
	defer list.mu.Unlock()
	if _, exist := list.devices[id]; exist {
		delete(list.devices, id)
		fmt.Printf("remove device: '%s' \n", id)
		return nil
	}
	return errors.New("device does not exist")
}

func (list *List) StartAliveWorker(publishRate time.Duration) func() {
	stop := make(chan struct{})
	go onAliveResponse(list)
	go checkDeviceNotResp(list)

	go func() {
		for {
			select {
			case <-stop:
				fmt.Println("Worker stopped")
				return
			default:
				token := client.Publish("device/pairing", 0, false, "pairing")
				token.Wait()
				if token.Error() != nil {
					fmt.Printf("Error publishing to topic: %v\n", token.Error())
				}
				time.Sleep(publishRate)
			}
		}
	}()

	return func() {
		close(stop)
	}
}

func (l *List) FindDevice(id string) (devices.IoTDevice, bool) {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// TODO: check why can't lock
	if device, have := l.devices[id]; have {
		return device, true
	}
	return nil, false
}

func (l *List) GetDevicesJSON() map[string]any {
	l.mu.RLock()
	defer l.mu.RUnlock()

	// Prepare a map for JSON response
	result := make(map[string]any)

	// Iterate over devices and extract their data
	for id, device := range l.devices {
		result[id] = device.GetDataResp() // Call GetData() to get Device struct
	}

	return result
}
func onAliveResponse(list *List) {
	client.Subscribe("device/paired", 0, func(c mqtt.Client, m mqtt.Message) {
		now := time.Now()
		deviceID := string(m.Payload())
		// type of device
		device, have := list.FindDevice(deviceID)
		if !have {

			newDevice := devices.NewDevice("PC", deviceID, deviceID, &client, func(device *devices.Device) {
				list.onStateChange(device, list)
			})

			if newDevice == nil {
				fmt.Println("Error creating new device")
				return
			}

			list.AddDevice(deviceID, newDevice)
			newDevice.SetLastCheck(&now)
			newDevice.ChangeState(devices.CONNECTED)
			//TODO: for do something before idle state

			newDevice.ChangeState(devices.IDLE)
			fmt.Println("device :", newDevice.GetData(), "is new connected")

			return
		}

		device.SetLastCheck(&now)
		if device.IsState(devices.DISCONNECTED) {
			device.ChangeState(devices.CONNECTED)
			//TODO: for do something before idle state

			device.ChangeState(devices.IDLE)
			fmt.Println("device :", device.GetData().Id, "is reconnected")
		}
	})
}

func checkDeviceNotResp(list *List) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		list.mu.RLock()
		devicesCopy := make([]devices.IoTDevice, 0, len(list.devices))
		for _, device := range list.devices {
			devicesCopy = append(devicesCopy, device) // Copy devices for safe access
		}
		list.mu.RUnlock()

		for _, device := range devicesCopy {
			if !device.IsState(devices.DISCONNECTED) && time.Since(*device.GetData().LastCheck) > (time.Second*10) {
				device.ChangeState(devices.DISCONNECTED)
				fmt.Println("device :", device.GetData().Id, "is disconnected")
			}
		}
	}
}
