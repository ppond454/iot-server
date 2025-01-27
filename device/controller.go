package device

import (
	"errors"
	"fmt"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type IController interface {
	NewController() (*List, error)
	CheckAliveWorker() func()
	FindDevice(device string)
}

type List struct {
	devices map[string]*Device
	mu      sync.Mutex
}

var client mqtt.Client = nil

func NewController(host string) (*List, error) {
	if client != nil {
		return nil, errors.New("controller already exists")
	}
	opts := mqtt.NewClientOptions().AddBroker(host)
	client = mqtt.NewClient(opts)
	return &List{devices: make(map[string]*Device)}, nil
}

func (list *List) Connect() error {
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("Error connecting to broker: &s", token.Error())
	}
	fmt.Printf("Connected to broker\n")
	return nil
}

func (list *List) AddDevice(id string, d *Device) (map[string]*Device, error) {
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
	if _, exist := list.devices[id]; exist {
		delete(list.devices, id)
		fmt.Printf("remove device: '%s' \n", id)
		return nil
	}
	return errors.New("device does not exist")
}

func (list *List) Disconnect() {
	client.Disconnect(250)
}

func (list *List) CheckAliveWorker(publishRate time.Duration) func() {
	stop := make(chan struct{})
	onAliveResponse(list)
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

func (list *List) FindDevice(id string) (*Device, bool) {
	list.mu.Lock()
	defer list.mu.Unlock()
	device, have := list.devices[id]
	if have {
		return device, true
	}
	return &Device{}, false
}

func onAliveResponse(list *List) {
	client.Subscribe("device/paired", 0, func(c mqtt.Client, m mqtt.Message) {
		now := time.Now()
		deviceID := string(m.Payload())
		device, have := list.FindDevice(deviceID)
		if !have {
			newDevice := NewDevice(deviceID, deviceID)
			list.AddDevice(deviceID, newDevice)
			newDevice.Connected(&now)

		} else {
			device.Connected(&now)
		}
	})
}

func checkDeviceNotResp(list *List) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for range ticker.C {
		list.mu.Lock()
		for _, device := range list.devices {
			if device.IsConnected() && time.Since(*device.GetData().lastCheck) > (time.Second*5) {
				device.Disconnect()
				fmt.Println("device :", device.Id, "is disconnected", device.isConnected)
			}
		}
		list.mu.Unlock()
	}
}
