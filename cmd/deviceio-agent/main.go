package main

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	_ "github.com/deviceio/agent/resources/filesystem"
	"github.com/deviceio/agent/transport"
	"github.com/deviceio/shared/logging"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

func main() {
	rand.Seed(time.Now().UnixNano()) //very important

	homedir, err := homedir.Dir()

	if err != nil {
		log.Fatal(err)
	}

	viper.SetConfigName("config")
	viper.AddConfigPath(fmt.Sprintf("%v/.deviceio/agent/", homedir))
	viper.AddConfigPath("/etc/deviceio/agent/")
	viper.AddConfigPath("/opt/deviceio/agent/")
	viper.AddConfigPath("c:/PROGRA~1/deviceio/agent/")
	viper.AddConfigPath("c:/ProgramData/deviceio/agent/")
	viper.AddConfigPath(".")

	viper.SetDefault("id", "")
	viper.SetDefault("tags", []string{})
	viper.SetDefault("transport.host", "127.0.0.1")
	viper.SetDefault("transport.port", 8975)
	viper.SetDefault("transport.allow_self_signed", "false")

	viper.BindEnv("id", "DEVICEIO_AGENT_ID")
	viper.BindEnv("tags", "DEVICEIO_AGENT_TAGS")
	viper.BindEnv("transport.host", "DEVICEIO_AGENT_TRANSPORT_HOST")
	viper.BindEnv("transport.port", "DEVICEIO_AGENT_TRANSPORT_PORT")
	viper.BindEnv("transport.allow_self_signed", "DEVICEIO_AGENT_TRANSPORT_INSECURE")

	transport.NewConnection(&logging.DefaultLogger{}).Dial(&transport.ConnectionOpts{
		ID:   viper.GetString("id"),
		Tags: viper.GetStringSlice("tags"),
		TransportAllowSelfSigned: viper.GetBool("transport.allow_self_signed"),
		TransportHost:            viper.GetString("transport.host"),
		TransportPort:            viper.GetInt("transport.port"),
	})

	<-make(chan bool)
}
