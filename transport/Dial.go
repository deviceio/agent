package transport

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/deviceio/shared/logging"

	"github.com/gorilla/websocket"
)

// Dial ...
func Dial(opts *Options) {
	dialer := &websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: opts.AllowTransportSelfSigned,
		},
	}

	url := fmt.Sprintf("wss://%v:%v/v1/connect", opts.TransportHost, opts.TransportPort)

	conn, _, err := dialer.Dial(url, http.Header{})

	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println(strings.Join(
		[]string{
			"Transport up:",
			fmt.Sprintf("LocalAddr=%v", conn.LocalAddr()),
			fmt.Sprintf("RemoteAddr=%v", conn.RemoteAddr()),
		},
		" ",
	))

	c := &connection{
		conn:    conn,
		options: opts,
		writech: make(chan []byte),
		logger:  &logging.DefaultLogger{},
	}

	go c.start()
}
