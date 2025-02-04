package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/ppond454/iot-backend/internal/model"
	mqttClient "github.com/ppond454/iot-backend/internal/mqtt"

	"net/http"

	io "github.com/googollee/go-socket.io"
)

func main() {

	client, disconnect, err := mqttClient.Connect("192.168.1.100:1883")
	defer disconnect()

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}

	server := io.NewServer(nil)
	defer server.Close()

	OnStateChange := func(device *model.Device, list *model.List) {
		if d, ok := list.FindDevice(device.Id); ok {
			fmt.Println("device update", d)
			server.BroadcastToNamespace("/", "device_update", Response{
				Success: true,
				Message: "device update",
				Data:    d.GetDataResp(),
			})
		}
	}

	params := &model.Params{
		Client:        client,
		OnStateChange: OnStateChange,
	}

	manager, err := model.New(params)

	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	manager.StartAliveWorker(time.Second)

	server.OnConnect("/", func(c io.Conn) error {
		c.SetContext(c.ID())
		return nil
	})

	server.OnEvent("/", "device_list", func(s io.Conn, msg string) {
		result := manager.GetDevicesJSON()
		s.Emit("device_list", Response{
			Success: true,
			Data:    result,
		})
	})

	server.OnEvent("/", "command", func(s io.Conn, body map[string]any) {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			s.Emit("command", Response{
				Success: false,
				Message: "Invalid JSON",
				Errors:  "Invalid JSON",
			})
			return
		}

		var input InputCommand
		if err := json.Unmarshal(jsonBody, &input); err != nil {
			s.Emit("command", Response{
				Success: false,
				Message: "Invalid Body",
				Errors:  "Invalid Body",
			})
			return
		}

		d, ok := manager.FindDevice(input.DeviceId)
		if !ok {
			s.Emit("command", Response{
				Success: false,
				Message: "Device not found",
				Errors:  "Device not found",
			})
			return
		}

		if !d.IsState(model.IDLE) {
			state := d.GetState()
			s.Emit("command", Response{
				Success: false,
				Message: fmt.Sprintf("device on state %s is unavailable", state.String()),
				Errors:  fmt.Sprintf("device %s is unavailable", input.DeviceId),
			})
			return
		}

		if err := d.RequestToggle(); err != nil {
			s.Emit("command", Response{
				Success: false,
				Message: err.Error(),
				Errors:  "request failed",
			})
			return
		}

		s.Emit("command", Response{
			Success: true,
			Data:    d.GetDataResp(),
		})
	})

	server.OnDisconnect("/", func(s io.Conn, reason string) {
		fmt.Println("closed =>", s.ID(), reason)
	})

	http.Handle("/socket.io/", server)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO error: %v", err)
		}
	}()

	PORT := ":8080"
	fmt.Printf("Server running on %s\n", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))

}

type Response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

type InputCommand struct {
	DeviceId string `json:"deviceId"`
}
