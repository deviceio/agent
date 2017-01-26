package transport

import (
	"io"
	"quantum/shared/protocol_v1"
)

// Options ...
type Options struct {
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
	HandleResource             func(*protocol_v1.Envelope, io.WriteCloser)
}
