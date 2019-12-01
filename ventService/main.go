package main

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	ventService "github.com/danilsh/simpleSmart/ventService/proto"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/micro/go-micro"
)

// VentServiceImpl structure used to store channel for stop signal
type VentServiceImpl struct {
	State            bool
	LastTurnOnTime   time.Time
	LastWorkDuration time.Duration
}

var ventServiceState = VentServiceImpl{State: false}

// TurnOff will turn off channel vent regardless of sensor values
func (state *VentServiceImpl) TurnOff(context.Context, *ventService.TurnOffMsg, *ventService.TurnOffResponse) error {
	return state.off()
}

func (state *VentServiceImpl) off() error {
	if state.State != false {
		state.State = false
		state.LastWorkDuration = time.Since(state.LastTurnOnTime)
	}

	if !mqttClient.IsConnectionOpen() {
		return errors.New("MQTT connection lost")
	}

	fmt.Println("Turn vent OFF")
	if token := mqttClient.Publish("ROOT/Actuators/BathVentActuator/State", 0, false, strconv.FormatBool(state.State)); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if token := mqttClient.Publish("ROOT/Actuators/BathVentActuator/LastWorkDuration", 0, false, state.LastWorkDuration.String()); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// TurnOn will turn on channel vent regardless of sensor values
func (state *VentServiceImpl) TurnOn(context.Context, *ventService.TurnOnMsg, *ventService.TurnOnResponse) error {
	return state.on()
}

func (state *VentServiceImpl) on() error {
	if state.State != true {
		state.State = true
		state.LastTurnOnTime = time.Now()
	}

	if !mqttClient.IsConnectionOpen() {
		return errors.New("MQTT connection lost")
	}

	fmt.Println("Turn vent ON")
	if token := mqttClient.Publish("ROOT/Actuators/BathVentActuator/State", 0, false, strconv.FormatBool(state.State)); token.Wait() && token.Error() != nil {
		return token.Error()
	}
	if token := mqttClient.Publish("ROOT/Actuators/BathVentActuator/LastTurnOnTime", 0, false, state.LastTurnOnTime.String()); token.Wait() && token.Error() != nil {
		return token.Error()
	}

	return nil
}

// ProcessSensorData will analyse received sensor data and manage channel vent state
func (state *VentServiceImpl) ProcessSensorData(humidity string) error {
	h, err := strconv.ParseFloat(strings.TrimSpace(humidity), 64)
	if err != nil {
		return err
	}

	// Для того, чтобы вентилятор не включался/выключался постоянно при дребезге сенсора
	// относительно пороговой точки, реализуем гистерезис
	if h > 66 {
		return state.on()
	}
	if h < 64 {
		return state.off()
	}

	return nil
}

var mqttClient mqtt.Client

func setupMqtt() bool {
	options := mqtt.NewClientOptions()
	// TODO: MQTT broker discovery
	options.AddBroker("192.168.1.22:1883")
	options.OnConnectionLost = func(client mqtt.Client, err error) {
		fmt.Println("MQTT connection lost")
	}
	options.OnReconnecting = func(clent mqtt.Client, opt *mqtt.ClientOptions) {
		fmt.Println("MQTT reconnecting")
	}
	options.OnConnect = func(client mqtt.Client) {
		client.SubscribeMultiple(
			map[string]byte{
				"ROOT/Sensors/DHT11_1/Temperature": 0,
				"ROOT/Sensors/DHT11_1/Humidity":    0,
				"ROOT/Sensors/DHT11_1/VCC":         0},
			func(client mqtt.Client, message mqtt.Message) {
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

func main() {
	if !setupMqtt() {
		return
	}
	defer mqttClient.Disconnect(0)

	srv := micro.NewService(micro.Name("ventService"))
	srv.Init()

	ventService.RegisterVentServiceHandler(srv.Server(), &ventServiceState)

	srv.Run()
}
