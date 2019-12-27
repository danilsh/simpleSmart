package main

import (
	"log"
	"strings"
	"time"
		
	mqtt "github.com/eclipse/paho.mqtt.golang"
	influx "github.com/influxdata/influxdb1-client/v2"
)

var mqttClient mqtt.Client

func setupMqtt() bool {
	options := mqtt.NewClientOptions()
	// TODO: MQTT broker discovery
	//options.AddBroker("192.168.1.187:1883")
	options.AddBroker("127.0.0.1:1883")
	options.OnConnectionLost = func(client mqtt.Client, err error) {
		log.Println("MQTT connection lost")
	}
	options.OnReconnecting = func(clent mqtt.Client, opt *mqtt.ClientOptions) {
		log.Println("MQTT reconnecting")
	}
	options.OnConnect = func(client mqtt.Client) {
		client.Subscribe("ROOT/Sensors/#", 0,
			func(client mqtt.Client, message mqtt.Message) {
				registerData(message.Topic(), string(message.Payload()));
				log.Print(strings.TrimSpace(string(message.Payload())))
				switch message.Topic() {
				case "ROOT/Sensors/DHT11_1/Temperature":
					log.Println("*C")
				case "ROOT/Sensors/DHT11_1/Humidity":
					log.Println("%")
				case "ROOT/Sensors/DHT11_1/VCC":
					log.Println("V")
				}
			})
	}
	mqttClient = mqtt.NewClient(options)

	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		log.Println("Error connecting to mqtt server")
		log.Println(token.Error())
		return false
	}

	return true
}

var influxDbClient influx.Client

func setupInflux() bool {
	client, err := influx.NewHTTPClient(influx.HTTPConfig { Addr: "http://127.0.0.1:8086", Timeout: 300 * time.Second })
	if err != nil {
		log.Println("Error connecting to InfluxDB server")
		log.Println(err.Error())
		return false
	}
	influxDbClient = client
	
	return true
}
