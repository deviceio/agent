package transport

import (
	"log"

	"github.com/deviceio/shared/logging"
	"github.com/deviceio/shared/protocol_v1"

	"github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
)

type connection struct {
	conn    *websocket.Conn
	options *Options
	writech chan []byte
	logger  logging.Logger
}

func (t *connection) start() {
	handshake := &protocol_v1.Handshake{
		AgentID:      t.options.ID,
		Hostname:     "TEST",
		Platform:     "TEST",
		Architecture: "TEST",
		Tags:         make([]string, 0),
	}

	handshakeb, err := proto.Marshal(handshake)

	if err != nil {
		log.Println("Error marshaling handshale", err)
	}

	envelope := &protocol_v1.Envelope{
		Type: protocol_v1.Envelope_Handshake,
		Data: handshakeb,
	}

	envelopeb, err := proto.Marshal(envelope)

	if err != nil {
		log.Println("Error marshalling envelope", err)
	}

	go t.readLoop()
	go t.writeLoop()

	t.writech <- envelopeb
}

func (t *connection) readLoop() {
	for {
		mt, mb, err := t.conn.ReadMessage()

		if err != nil {
			log.Println(err)
			t.conn.Close()
			return
		}

		if mt != websocket.BinaryMessage {
			t.logger.Error("Invalid message type")
			t.conn.Close()
			return
		}

		envelope := &protocol_v1.Envelope{}

		if err := proto.Unmarshal(mb, envelope); err != nil {
			t.logger.Error(err.Error())
			t.conn.Close()
			return
		}

		w := &writer{
			resv:  make(chan []byte),
			close: make(chan bool),
		}

		go func() {
			t.options.HandleResource(envelope, w)
		}()

		go func(w *writer) {
			for {
				select {
				case b := <-w.resv:
					t.writech <- b
				case <-w.close:
					return
				}
			}
		}(w)
	}
}

func (t *connection) writeLoop() {
	for {
		err := t.conn.WriteMessage(websocket.BinaryMessage, <-t.writech)
		if err != nil {
			log.Println(err)
			t.conn.Close()
			return
		}
	}
}
