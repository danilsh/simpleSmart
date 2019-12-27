package main

import (
	"github.com/micro/go-micro"
)
	
func main() {
	if !setupMqtt() {
		return
	}
	defer mqttClient.Disconnect(0)
	
	if !setupInflux() {
		return
	}
	defer influxDbClient.Close()
	
	srv := micro.NewService(micro.Name("historyCaptureService"))
	srv.Init()

	srv.Run()	
}
