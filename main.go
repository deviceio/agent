package main

import (
	"log"
	"math/rand"
	"time"

	_ "github.com/deviceio/agent/resources/filesystem"
	"github.com/deviceio/agent/transport"
	"github.com/deviceio/shared/config"
)

func main() {
	rand.Seed(time.Now().Unix()) //very important

	var configuration struct {
		DisableTransportKeyPinning bool
		AllowTransportSelfSigned   bool
		ID                         string
		TransportHost              string
		TransportPort              int
		PasscodeHash               string
		PasscodeSalt               string
		Tags                       []string
	}

	config.SetConfigStruct(&configuration)
	config.AddFileName("config.json")
	config.AddFilePath("/etc/deviceio/agent")
	config.AddFilePath("c:/ProgramData/deviceio/agent")

	if err := config.Parse(); err != nil {
		log.Fatal(err)
	}

	(&transport.Connection{}).Dial(&transport.ConnectionOpts{
		DisableTransportKeyPinning: configuration.DisableTransportKeyPinning,
		AllowTransportSelfSigned:   configuration.AllowTransportSelfSigned,
		ID:            configuration.ID,
		TransportHost: configuration.TransportHost,
		TransportPort: configuration.TransportPort,
		PasscodeHash:  configuration.PasscodeHash,
		PasscodeSalt:  configuration.PasscodeSalt,
		Tags:          configuration.Tags,
	})

	<-make(chan bool)
}
