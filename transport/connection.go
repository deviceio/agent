package transport

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"

	"time"

	"github.com/deviceio/shared/logging"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
	"github.com/jpillora/backoff"
)

// Connection represents our upstream connection to a hub.
type Connection struct {
	opts      *ConnectionOpts
	reconnect int
	jitter    int
	backoff   *backoff.Backoff
	logger    logging.Logger
}

// NewConnection creates a new instance of the Connection type
func NewConnection(logger logging.Logger) *Connection {
	return &Connection{
		logger: logger,
	}
}

// Dial attempts to connect to the upstream hub. If dialing fails a backoff
// algorithm is applied during reconnection attempts to alleviate load on a hub
// that disappears momentarily.
func (t *Connection) Dial(opts *ConnectionOpts) {
	t.opts = opts

	t.backoff = &backoff.Backoff{
		Max:    5 * time.Second,
		Jitter: true,
	}

	for {
		err := t.run()

		t.logger.Warn(fmt.Sprintf(
			"Transport Failure %v:%v : %v",
			t.opts.TransportHost,
			t.opts.TransportPort,
			err.Error(),
		))

		wait := t.backoff.Duration()

		if wait >= t.backoff.Max {
			t.backoff.Reset()
		}

		t.logger.Info(fmt.Sprintf(
			"Reconnect %v:%v in %v seconds",
			t.opts.TransportHost,
			t.opts.TransportPort,
			wait,
		))

		time.Sleep(wait)
	}
}

// run conducts the setup of the multiplexed tcp stream server to the hub and registers
// the base http server to be served over the multiplexed connection.
func (t *Connection) run() error {
	dialer := &websocket.Dialer{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: t.opts.TransportAllowSelfSigned,
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

	t.logger.Info(strings.Join(
		[]string{
			"Transport up:",
			fmt.Sprintf("LocalAddr=%v", conn.LocalAddr()),
			fmt.Sprintf("RemoteAddr=%v", conn.RemoteAddr()),
		},
		" ",
	))

	server, _ := yamux.Server(conn.UnderlyingConn(), nil)

	Router.HandleFunc("/info", t.httpGetInfo)

	return http.Serve(server, Router)
}

// httpGetInfo provides basic information about this device a hub needs to properly
// manage the transport connection.
func (t *Connection) httpGetInfo(w http.ResponseWriter, r *http.Request) {
	type info struct {
		ID           string
		Hostname     string
		Architecture string
		Platform     string
		Tags         []string
	}

	hostname, err := os.Hostname()

	if err != nil {
		hostname = "Unknown"
	}

	err = json.NewEncoder(w).Encode(&info{
		ID:           t.opts.ID,
		Tags:         t.opts.Tags,
		Hostname:     hostname,
		Architecture: runtime.GOARCH,
		Platform:     runtime.GOOS,
	})

	if err != nil {
		t.logger.Error(err.Error())
		w.WriteHeader(500)
		w.Write([]byte(""))
		return
	}
}
