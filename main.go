package main

import (
	"log"
	"math/rand"
	"time"

	_ "github.com/deviceio/agent/resources/filesystem"
	"github.com/deviceio/agent/transport"
	"github.com/deviceio/shared/config"
	"github.com/deviceio/shared/logging"
)

func main() {
	rand.Seed(time.Now().Unix()) //very important

	var configuration struct {
		ID         string
		Tags       []string
		Transports []struct {
			Host              string
			Port              int
			PasscodeHash      string
			PasscodeSalt      string
			DisableKeyPinning bool
			AllowSelfSigned   bool
		}
	}

	config.SetConfigStruct(&configuration)
	config.AddFileName("config.json")
	config.AddFilePath("/etc/deviceio/agent")
	config.AddFilePath("c:/ProgramData/deviceio/agent")

	if err := config.Parse(); err != nil {
		log.Fatal(err)
	}

	for _, transportConfig := range configuration.Transports {
		c := transport.NewConnection(&logging.DefaultLogger{})

		c.Dial(&transport.ConnectionOpts{
			ID:   configuration.ID,
			Tags: configuration.Tags,
			TransportAllowSelfSigned:   transportConfig.AllowSelfSigned,
			TransportDisableKeyPinning: transportConfig.DisableKeyPinning,
			TransportHost:              transportConfig.Host,
			TransportPasscodeHash:      transportConfig.PasscodeHash,
			TransportPasscodeSalt:      transportConfig.PasscodeSalt,
			TransportPort:              transportConfig.Port,
		})
	}

	<-make(chan bool)
}
