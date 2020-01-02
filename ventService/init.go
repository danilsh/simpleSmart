package main

import (
	"log"
	"time"
		
	mqtt "github.com/eclipse/paho.mqtt.golang"
)

var mqttClient mqtt.Client

func setupMqtt() bool {
	options := mqtt.NewClientOptions()
	// TODO: MQTT broker discovery
	options.AddBroker("127.0.0.1:1883")
	options.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Println("MQTT connection lost")
	}
	options.OnReconnecting = func(clent mqtt.Client, opt *mqtt.ClientOptions) {
		log.Println("MQTT reconnecting")
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
				//log.Print(strings.TrimSpace(string(message.Payload())))
				//switch message.Topic() {
				//case "ROOT/Sensors/DHT11_1/Temperature":
				//	log.Println("*C")
				//case "ROOT/Sensors/DHT11_1/Humidity":
				//	log.Println("%")
				//case "ROOT/Sensors/DHT11_1/VCC":
				//	log.Println("V")
				//}
			})
			log.Println("MQTT connected");
	}
	mqttClient = mqtt.NewClient(options)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Error connecting to mqtt server")
		log.Println(token.Error())
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

func initCalc() {
	calcTick = time.NewTicker(calcTimeout)
	go calcLoop()
}

