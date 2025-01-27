package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/ppond454/iot-backend/device"
)

func main() {
	controller, err := device.NewController("192.168.1.100:1883")

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	if err := controller.Connect(); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	controller.CheckAliveWorker(time.Second * 2)
	test()

	select {}

}

func test() {
	opts := mqtt.NewClientOptions().AddBroker("192.168.1.100:1883")
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return

	}

	client.Subscribe("device/pairing", 0, func(c mqtt.Client, m mqtt.Message) {
		// fmt.Printf("Received message: [%s] %s\n", m.Topic(), m.Payload())
		go func() {
			for i := 0; i < 2; i++ {
				token := client.Publish("device/paired", 0, false, fmt.Sprintf("device%d", i))
				token.Wait()
				if token.Error() != nil {
					fmt.Printf("Error publishing to topic: %v\n", token.Error())
				}

			}
		}()
		time.Sleep(time.Second * 3)
		client.Unsubscribe("device/pairing")
		time.Sleep(time.Second * 3)

		test()
	})

}
