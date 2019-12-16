package main

import (
	"strconv"
	"strings"
	"fmt"
	"time"
)

// Будем собирать данные в два массива:
// В первом массиве будем хранить информацию с датчиков из "сырых" помещений
// Во втором массиве будем хранить информацию с датчиков из "сухих" помещений
//
// Данные будем собирать по мере поступления, т.е. с той скоростью, с которой их отправляют датчики.
// При этом, нам важно получать информацию хотя бы от одного датчика из "влажного" помещения.
//
// Будем контролировать то, что датчики "живы" с помощью таймера. При получении значения
// хотя бы от одного датчика из группы "сухие" или "влажные" будем сбрасывать соответствующий
// таймер. Если какой-то из таймеров достигнет тайм-аута, то будем считать все датчики соответствующей
// группы "мёртвыми". Это состояние группы так же сохраняется. Оно будет сброшено, как только
// в соответствующей группе будет получено хотя бы одно значение.

type DHTSensorData struct {
	Temperature float64
	Humidity float64
	VCC float64
}

// TODO: как вариант, можно переделать систему передачи информации от модуля регистрации в расчётный
// модуль. Информацию можно передавать через каналы. Но сейчас это не совсем удобно, т.к. нужно выбирать
// максимальные-минимальные сигналы от датчиков и т.п.
var wetSensors = make([]DHTSensorData, 2)
var drySensors = make([]DHTSensorData, 2)
var wetSensorsAlive bool = false
var drySensorsAlive bool = false
const sensorTimeout = time.Second * 90
var wetTimer *time.Timer
var dryTimer *time.Timer

func registerData(topic string, message string) {
	value, err := strconv.ParseFloat(strings.TrimSpace(message), 64)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// TODO: Нужно с этим ужасом что-то делать.
	// А если датчиков будет сто? (не будет)
	// Некрасиво, кароч
	switch topic {
	case "ROOT/Sensors/DHT11_1/Temperature":
		wetSensors[0].Temperature = value
		wetSensorsAlive = true
		wetTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/DHT11_1/Humidity":
		wetSensors[0].Humidity = value
		wetSensorsAlive = true
		wetTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/DHT11_1/VCC":
		wetSensors[0].VCC = value
		wetSensorsAlive = true
		wetTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/WetRooms/2/Temperature":
		wetSensors[1].Temperature = value
		wetSensorsAlive = true
		wetTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/WetRooms/2/Humidity":
		wetSensors[1].Humidity = value
		wetSensorsAlive = true
		wetTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/WetRooms/2/VCC":
		wetSensors[1].VCC = value
		wetSensorsAlive = true
		wetTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/DryRooms/1/Temperature":
		drySensors[0].Temperature = value
		drySensorsAlive = true
		dryTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/DryRooms/1/Humidity":
		drySensors[0].Humidity = value
		drySensorsAlive = true
		dryTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/DryRooms/1/VCC":
		drySensors[0].VCC = value
		drySensorsAlive = true
		dryTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/DryRooms/2/Temperature":
		drySensors[1].Temperature = value
		drySensorsAlive = true
		dryTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/DryRooms/2/Humidity":
		drySensors[1].Humidity = value
		drySensorsAlive = true
		dryTimer.Reset(sensorTimeout)
	case "ROOT/Sensors/DryRooms/2/VCC":
		drySensors[1].VCC = value
		drySensorsAlive = true
		dryTimer.Reset(sensorTimeout)
	}
}

func wetTimerControl() {
	for {
		<-wetTimer.C
		wetSensorsAlive = false
		fmt.Println("Wet sensors dead")
		wetTimer.Reset(sensorTimeout)
	}
}

func dryTimerControl() {
	for {
		<-dryTimer.C
		drySensorsAlive = false
		fmt.Println("Dry sensors dead")
		dryTimer.Reset(sensorTimeout)
	}
}
