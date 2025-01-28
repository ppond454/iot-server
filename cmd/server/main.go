package main

import (
	"fmt"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	manager "github.com/ppond454/iot-backend/internal/manager"
	mqttClient "github.com/ppond454/iot-backend/pkg/mqtt"
)

func main() {

	client, _, err := mqttClient.Connect("192.168.1.100:1883")

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	manager, err := manager.New(client)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	manager.StartAliveWorker(time.Second)

	test()

	select {}

}

func test() {
	client, disconnect, err := mqttClient.Connect("192.168.1.100:1883")

	if err != nil {
		fmt.Println("Error: ", err)
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
		disconnect()
		time.Sleep(time.Second * 15)
		test()
	})

}
