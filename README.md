# **IoT Device Management System with MQTT & Socket.IO**

## **Project Overview**

This project is an **IoT Device Management System** that enables real-time communication between IoT devices and a central server using **MQTT** and **Socket.IO**. It provides a WebSocket-based API for **client applications (mobile/web)** to interact with IoT devices, allowing users to **monitor and control** them in real time.

---

## **Features**

✅ **Device Connection Management** – Track device connections and disconnections.  
✅ **Real-time Communication** – Send and receive device updates via **Socket.IO**.  
✅ **MQTT Integration** – Devices publish and subscribe to topics for data exchange.  
✅ **Device Status Updates** – Automatically update device states (e.g., ON/OFF).  
✅ **WebSocket API for Clients** – Allow web and mobile apps to interact with devices.  
❌ (TODO: implement) **Authentication Middleware**

---

## **Technology Stack**

- **Golang** – High-performance backend development
- **Socket.IO** – WebSocket-based real-time communication
- **MQTT (Eclipse Mosquitto)** – Device communication protocol

---

## **Installation & Setup**

### **1. Clone the Repository**

```bash
git clone https://github.com/ppond454/iot-server.git
cd iot-server
```

### **2. Install Dependencies**

```bash
go mod tidy
```

### **3. Run 🏃‍♂️**

```bash
go run main.go
```
