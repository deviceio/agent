package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime"
	"time"

	"golang.org/x/sys/windows/svc"
	"golang.org/x/sys/windows/svc/debug"

	"github.com/alecthomas/kingpin"
	"github.com/deviceio/agent/installation"
	_ "github.com/deviceio/agent/resources/filesystem"
	"github.com/deviceio/agent/transport"
	"github.com/deviceio/shared/config"
	"github.com/deviceio/shared/logging"
)

var (
	cli = kingpin.New("cli", "Deviceio Agent Command Line Interface")

	installCommand           = cli.Command("install", "installs the agent to the system")
	installOrganization      = installCommand.Flag("org", "Specify the organization responsible for this agent installation").Short('o').Required().String()
	installTransportHost     = installCommand.Flag("transport-host", "The hostname or ip address of your hub installation").Default("127.0.0.1").Short('h').String()
	installTransportPort     = installCommand.Flag("transport-port", "The port number of your hub's gateway binding").Default("8975").Short('p').Int()
	installTransportInsecure = installCommand.Flag("transport-insecure", "Do not check validity of the hub TLS Certificate").Default("false").Short('i').Bool()

	startCommand = cli.Command("start", "Starts the agent and connects to the hub transport")
	startConfig  = startCommand.Arg("config", "The configuration file to load").Required().String()

	serviceCommand = cli.Command("service", "Runs the agent as a service on compatable systems")
	serviceConfig  = serviceCommand.Arg("config", "The configuration file to load").Required().String()
)

func main() {
	rand.Seed(time.Now().UnixNano()) //very important

	switch kingpin.MustParse(cli.Parse(os.Args[1:])) {
	case installCommand.FullCommand():
		if err := installation.Install(
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
	case serviceCommand.FullCommand():
		if runtime.GOOS == "windows" {
			svc.Run("deviceio agent", &winsvc{})
			config.AddConfigPath(*serviceConfig)
		} else {
			panic("service mode is not supported on this system")
		}
	default:
		return
	}

	var configuration *installation.Config
	config.SetConfigStruct(&configuration)
	config.AddConfigPath(*startConfig)

	if err := config.Parse(); err != nil {
		log.Fatal(err)
	}

	c := transport.NewConnection(&logging.DefaultLogger{})

	c.Dial(&transport.ConnectionOpts{
		ID:   configuration.ID,
		Tags: configuration.Tags,
		TransportAllowSelfSigned: configuration.TransportAllowSelfSigned,
		TransportHost:            configuration.TransportHost,
		TransportPort:            configuration.TransportPort,
	})

	<-make(chan bool)
}

var elog debug.Log

type winsvc struct{}

func (m *winsvc) Execute(args []string, r <-chan svc.ChangeRequest, changes chan<- svc.Status) (ssec bool, errno uint32) {
	const cmdsAccepted = svc.AcceptStop | svc.AcceptShutdown | svc.AcceptPauseAndContinue
	changes <- svc.Status{State: svc.StartPending}
	changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
loop:
	for {
		select {
		case c := <-r:
			switch c.Cmd {
			case svc.Interrogate:
				changes <- c.CurrentStatus
			case svc.Stop, svc.Shutdown:
				break loop
			case svc.Pause:
				changes <- svc.Status{State: svc.Paused, Accepts: cmdsAccepted}
			case svc.Continue:
				changes <- svc.Status{State: svc.Running, Accepts: cmdsAccepted}
			default:
				elog.Error(1, fmt.Sprintf("unexpected control request #%d", c))
			}
		}
	}
	changes <- svc.Status{State: svc.StopPending}
	return
}
