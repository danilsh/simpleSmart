package main

import (
	"fmt"
	"strings"
	"time"
		
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

func setupMqtt() bool {
	options := mqtt.NewClientOptions()
	// TODO: MQTT broker discovery
	options.AddBroker("192.168.1.187:1883")
	options.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Println("MQTT connection lost")
	}
	options.OnReconnecting = func(clent mqtt.Client, opt *mqtt.ClientOptions) {
		fmt.Println("MQTT reconnecting")
	}
	options.OnConnect = func(client mqtt.Client) {
		client.Subscribe("ROOT/Sensors/#", 0,
//		client.SubscribeMultiple(
//			map[string]byte{
//				"ROOT/Sensors/DHT11_1/Temperature": 0,
//				"ROOT/Sensors/DHT11_1/Humidity":    0,
//				"ROOT/Sensors/DHT11_1/VCC":         0},
			func(client mqtt.Client, message mqtt.Message) {
				registerData(message.Topic(), string(message.Payload()));
				fmt.Print(strings.TrimSpace(string(message.Payload())))
				switch message.Topic() {
				case "ROOT/Sensors/DHT11_1/Temperature":
					fmt.Println("*C")
				case "ROOT/Sensors/DHT11_1/Humidity":
					fmt.Println("%")
					if err := ventServiceState.ProcessSensorData(string(message.Payload())); err != nil {
						fmt.Println(err.Error())
					}
				case "ROOT/Sensors/DHT11_1/VCC":
					fmt.Println("V")
				}
			})
	}
	mqttClient = mqtt.NewClient(options)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		fmt.Println("Error connecting to mqtt server")
		fmt.Println(token.Error())
		return false
	}

	return true
}

func initTimers() {
	wetTimer = time.NewTimer(sensorTimeout)
	go wetTimerControl()
	dryTimer = time.NewTimer(sensorTimeout)
	go dryTimerControl()
}

