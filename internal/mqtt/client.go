package mqtt

import (
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

func Connect(host string) (mqtt.Client, func(), error) {
	opts := mqtt.NewClientOptions().AddBroker(host)
	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return nil, nil, fmt.Errorf("Error connecting to broker: &s", token.Error())
	}
	fmt.Printf("Connected to broker\n")

	return client, func() {
		disconnect(client)
	}, nil
}

func disconnect(client mqtt.Client) {
	client.Disconnect(250)
}
