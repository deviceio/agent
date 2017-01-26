package main

import (
	"io"
	"log"
	"quantum/agent/resources"
	"quantum/agent/transport"
	"quantum/shared/config"
	"quantum/shared/protocol_v1"
)

func main() {
	var configuration struct {
		DisableTransportKeyPinning bool
		AllowTransportSelfSigned   bool
		ID                         string
		TransportHost              string
		TransportPort              int
		PasscodeHash               string
		PasscodeSalt               string
		ReconnectInterval          int
		ReconnectJitter            int
		Tags                       []string
	}

	config.SetConfigStruct(&configuration)
	config.AddFileName("config.json")
	config.AddFilePath("/etc/quantum/agent")
	config.AddFilePath("c:/ProgramData/quantum/agent")

	if err := config.Parse(); err != nil {
		log.Fatal(err)
	}

	reshandler := resources.NewHandler()

	go transport.Dial(&transport.Options{
		DisableTransportKeyPinning: configuration.DisableTransportKeyPinning,
		AllowTransportSelfSigned:   configuration.AllowTransportSelfSigned,
		ID:                configuration.ID,
		TransportHost:     configuration.TransportHost,
		TransportPort:     configuration.TransportPort,
		PasscodeHash:      configuration.PasscodeHash,
		PasscodeSalt:      configuration.PasscodeSalt,
		ReconnectInterval: configuration.ReconnectInterval,
		ReconnectJitter:   configuration.ReconnectJitter,
		Tags:              configuration.Tags,
		HandleResource: func(env *protocol_v1.Envelope, w io.WriteCloser) {
			reshandler.Handle(env, w)
		},
	})

	<-make(chan bool)
}
