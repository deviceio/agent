package main

import (
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/alecthomas/kingpin"
	"github.com/deviceio/agent/installer"
	_ "github.com/deviceio/agent/resources/filesystem"
	"github.com/deviceio/agent/transport"
	"github.com/deviceio/shared/config"
	"github.com/deviceio/shared/logging"
)

var (
	cli = kingpin.New("cli", "Deviceio Agent Command Line Interface")

	installCommand           = cli.Command("install", "installs the agent to the system")
	installOrganization      = installCommand.Flag("org", "Specify the organization responsible for this agent installation").Short('o').Required().String()
	installTransportHost     = installCommand.Flag("transport-host", "The hostname or ip address of your hub installation").Short('h').Required().String()
	installTransportPort     = installCommand.Flag("transport-port", "The port number of your hub's gateway binding").Short('p').Required().Int()
	installTransportInsecure = installCommand.Flag("transport-insecure", "Do not check validity of the hub TLS Certificate").Short('i').Required().Bool()

	startCommand = cli.Command("start", "Starts the agent and connects to the hub transport")
	startConfig  = startCommand.Arg("config", "The configuration file to load").Required().String()
)

func main() {
	rand.Seed(time.Now().UnixNano()) //very important

	switch kingpin.MustParse(cli.Parse(os.Args[1:])) {
	case installCommand.FullCommand():
		if err := installer.Install(
			*installOrganization,
			*installTransportHost,
			*installTransportPort,
			*installTransportInsecure,
		); err != nil {
			panic(err)
		}

		return
	case startCommand.FullCommand():
		break
	}

	var configuration struct {
		ID                         string
		Tags                       []string
		TransportHost              string
		TransportPort              int
		TransportPasscodeHash      string
		TransportPasscodeSalt      string
		TransportDisableKeyPinning bool
		TransportAllowSelfSigned   bool
	}

	config.SetConfigStruct(&configuration)
	config.AddConfigPath(*startConfig)

	if err := config.Parse(); err != nil {
		log.Fatal(err)
	}

	c := transport.NewConnection(&logging.DefaultLogger{})

	c.Dial(&transport.ConnectionOpts{
		ID:   configuration.ID,
		Tags: configuration.Tags,
		TransportAllowSelfSigned:   configuration.TransportAllowSelfSigned,
		TransportDisableKeyPinning: configuration.TransportDisableKeyPinning,
		TransportHost:              configuration.TransportHost,
		TransportPasscodeHash:      configuration.TransportPasscodeHash,
		TransportPasscodeSalt:      configuration.TransportPasscodeSalt,
		TransportPort:              configuration.TransportPort,
	})

	<-make(chan bool)
}
