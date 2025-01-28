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
	New(mqtt.Client) (*List, error)
	StartAliveWorker() func()
	FindDevice(device string)
}

type List struct {
	devices map[string]devices.IoTDevice
	mu      sync.Mutex
}

var client mqtt.Client = nil

func New(_client mqtt.Client) (*List, error) {
	if client != nil {
		return nil, errors.New("controller already exists")
	}
	client = _client
	return &List{devices: make(map[string]devices.IoTDevice)}, nil
}

func (list *List) AddDevice(id string, d devices.IoTDevice) (map[string]devices.IoTDevice, error) {
	list.mu.Lock()
	defer list.mu.Unlock()
	if _, exist := list.devices[id]; exist {
		fmt.Printf("Device '%s' already exists\n", id)
		return nil, fmt.Errorf("Device '%s' already exists", id)
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

func (list *List) FindDevice(id string) (devices.IoTDevice, bool) {
	list.mu.Lock()
	defer list.mu.Unlock()
	device, have := list.devices[id]
	if have {
		return device, true
	}
	return nil, false
}

func onAliveResponse(list *List) {
	client.Subscribe("device/paired", 0, func(c mqtt.Client, m mqtt.Message) {
		now := time.Now()
		deviceID := string(m.Payload())
		// type of device
		device, have := list.FindDevice(deviceID)
		if !have {
			newDevice := devices.NewDevice("PC", deviceID, deviceID)

			if newDevice == nil {
				fmt.Println("Error creating new device")
				return
			}

			list.AddDevice(deviceID, newDevice)
			newDevice.Connected(&now)
			// newDevice.
			fmt.Println("device :", newDevice.GetData(), "is new connected")

			return
		}

		if !device.IsConnected() {
			fmt.Println("device :", device.GetData().Id, "is reconnected")
		}
		device.Connected(&now)

		for _, l := range list.devices {
			fmt.Println("device :", l.GetData())
		}
	})
}

func checkDeviceNotResp(list *List) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		list.mu.Lock()
		for _, device := range list.devices {
			if device.IsConnected() && time.Since(*device.GetData().LastCheck) > (time.Second*10) {
				device.Disconnect()
				fmt.Println("device :", device.GetData().Id, "is disconnected")
			}
		}
		list.mu.Unlock()
	}
}
