package transport

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"time"

	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
	"github.com/jpillora/backoff"
)

type Connection struct {
	opts      *ConnectionOpts
	reconnect int
	jitter    int
	backoff   *backoff.Backoff
}

func (t *Connection) Dial(opts *ConnectionOpts) {
	t.opts = opts

	t.backoff = &backoff.Backoff{
		Max:    5 * time.Second,
		Jitter: true,
	}

	for {
		err := t.run()
		log.Println("Transport failure:", err)
		wait := t.backoff.Duration()

		if wait >= t.backoff.Max {
			t.backoff.Reset()
		}

		log.Println(fmt.Sprintf("Reconnect in %v seconds", wait))
		time.Sleep(wait)
	}
}

func (t *Connection) run() error {
	dialer := &websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: t.opts.AllowTransportSelfSigned,
		},
	}

	url := fmt.Sprintf(
		"wss://%v:%v/v1/connect",
		t.opts.TransportHost,
		t.opts.TransportPort,
	)

	conn, _, err := dialer.Dial(url, http.Header{})

	if err != nil {
		return err
	}

	log.Println(strings.Join(
		[]string{
			"Transport up:",
			fmt.Sprintf("LocalAddr=%v", conn.LocalAddr()),
			fmt.Sprintf("RemoteAddr=%v", conn.RemoteAddr()),
		},
		" ",
	))

	server, _ := yamux.Server(conn.UnderlyingConn(), nil)

	Router.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		encd := json.NewEncoder(w)
		encd.Encode(t.opts)
	})

	t.backoff.Reset()

	return http.Serve(server, Router)
}
