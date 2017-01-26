package types

import (
	"fmt"
	"io"
	"os"

	"github.com/deviceio/shared/protocol_v1"

	"github.com/golang/protobuf/proto"
)

// Core_FileSystem ...
type Core_FileSystem struct{}

// Call ...
func (t *Core_FileSystem) Call(member string, args map[string][]byte, env *protocol_v1.Envelope, w io.WriteCloser) {
	switch member {
	case "read":
		t.read(args, env, w)
	default:
		t.writeerr(fmt.Sprintf("No such member '%v'", member), env, w)
	}
}

func (t *Core_FileSystem) read(args map[string][]byte, env *protocol_v1.Envelope, w io.WriteCloser) {
	path, ok := args["path"]

	if !ok {
		t.writeerr("path argument must be supplied", env, w)
		return
	}

	pathstr := string(path)

	_, err := os.Stat(pathstr)

	if os.IsNotExist(err) {
		t.writeerr(fmt.Sprintf("path does not exist"), env, w)
		return
	} else if err != nil {
		t.writeerr(fmt.Sprintf("Stat error on path: %v", err), env, w)
		return
	}

	fh, err := os.Open(pathstr)

	if err != nil {
		t.writeerr(fmt.Sprintf("Error opening file: %v", err), env, w)
		return
	}

	fbytes := make([]byte, 250)

	for {
		n, err := fh.Read(fbytes)

		if err == io.EOF {
			return
		} else if err != nil {
			t.writeerr(fmt.Sprintf("Error reading file: %v", err), env, w)
			return
		}

		t.writebytes(fbytes[:n], env, w)
	}
}

func (t *Core_FileSystem) writebytes(b []byte, env *protocol_v1.Envelope, w io.WriteCloser) {
	bm := &protocol_v1.Bytes{
		Value: b,
	}

	msgb, _ := proto.Marshal(bm)

	env.Type = protocol_v1.Envelope_Bytes
	env.Data = msgb

	envb, _ := proto.Marshal(env)

	w.Write(envb)
}

func (t *Core_FileSystem) writeerr(err string, env *protocol_v1.Envelope, w io.WriteCloser) {
	e := &protocol_v1.Error{
		Message: err,
	}

	eb, _ := proto.Marshal(e)

	env.Type = protocol_v1.Envelope_Error
	env.Data = eb

	envb, _ := proto.Marshal(env)

	w.Write(envb)
}
