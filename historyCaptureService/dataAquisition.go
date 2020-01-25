package main

import (
	"strconv"
	"strings"
	"log"
	"time"
	
	influx "github.com/influxdata/influxdb1-client/v2"
)

func registerData(topic string, message string) {
	value, err := strconv.ParseFloat(strings.TrimSpace(message), 64)
	if err != nil {
		log.Println(err.Error())
		return
	}
	var fields map[string]interface{}
	// TODO: Нужно с этим ужасом что-то делать.
	// А если датчиков будет сто? (не будет)
	// Некрасиво, кароч
	switch topic {
	case "ROOT/Sensors/WetRooms/1/Temperature":
		fields = map[string]interface{}{ "Temperature" : value }
	case "ROOT/Sensors/WetRooms/1/Humidity":
		fields = map[string]interface{}{ "Humidity" : value }
	case "ROOT/Sensors/WetRooms/1/VCC":
		fields = map[string]interface{}{ "VCC" : value }
	case "ROOT/Sensors/WetRooms/2/Temperature":
		fields = map[string]interface{}{ "Temperature" : value }
	case "ROOT/Sensors/WetRooms/2/Humidity":
		fields = map[string]interface{}{ "Humidity" : value }
	case "ROOT/Sensors/WetRooms/2/VCC":
		fields = map[string]interface{}{ "VCC" : value }
	case "ROOT/Sensors/DryRooms/1/Temperature":
		fields = map[string]interface{}{ "Temperature" : value }
	case "ROOT/Sensors/DryRooms/1/Humidity":
		fields = map[string]interface{}{ "Humidity" : value }
	case "ROOT/Sensors/DryRooms/1/VCC":
		fields = map[string]interface{}{ "VCC" : value }
	case "ROOT/Sensors/DryRooms/2/Temperature":
		fields = map[string]interface{}{ "Temperature" : value }
	case "ROOT/Sensors/DryRooms/2/Humidity":
		fields = map[string]interface{}{ "Humidity" : value }
	case "ROOT/Sensors/DryRooms/2/VCC":
		fields = map[string]interface{}{ "VCC" : value }
	}
	
	point,err := influx.NewPoint (
		"sensors",
		map[string]string {
			"topic": topic,
		},
		fields,
		time.Now())
	if err != nil {
		log.Println("NewPoint() error: ")
		log.Println(err.Error())
		return
	}
	
	batchpoints, err := influx.NewBatchPoints(influx.BatchPointsConfig { Database: "simplesmart", RetentionPolicy: "sensors_retention" })
	if err != nil {
		log.Println(err.Error())
		return
	}
	batchpoints.AddPoint(point)
	err = influxDbClient.Write(batchpoints)
	if err != nil {
		log.Println("Write() error: ")
		log.Println(err.Error())
		return
	}
}
