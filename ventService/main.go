package main

import (
	ventService_RPC "ventService/proto"

	"github.com/micro/go-micro"
)

func main() {
	if !setupMqtt() {
		return
	}
	defer mqttClient.Disconnect(0)
	
	initTimers()
	initCalc()

	srv := micro.NewService(micro.Name("ventService"))
	srv.Init()

	ventService_RPC.RegisterVentService_RPCHandler(srv.Server(), &ventServiceState)

	srv.Run()
}
