package main

import (
	"fmt"

	"github.com/ppond454/iot-backend/device"
)

func main() {
	controllers := device.NewController()

	device1 := device.New("device1", "Smart Switch")
	controllers.AddDevice("device1", device1)

	device2 := device.New("device2", "Smart Thermostat")
	controllers.AddDevice("device2", device2)

	pingResult, err := device1.Ping()
	if err != nil {
		fmt.Printf("Ping failed for device1: %v\n", err)
	} else {
		fmt.Printf("Ping result for device1: %.2f\n", pingResult)
	}

	isConnected := device1.IsConnected()
	fmt.Printf("Device1 is connected: %v\n", isConnected)

	deviceData := device1.GetData()
	fmt.Printf("Device1 data: %+v\n", deviceData)
}
