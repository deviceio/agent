package resources

import (
	"io"
	"quantum/shared/protocol_v1"
)

// Resource ...
type Resource interface {
	Call(string, map[string][]byte, *protocol_v1.Envelope, io.WriteCloser)
}
