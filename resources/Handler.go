package resources

import (
	"io"
	"quantum/agent/resources/types"
	"quantum/shared/logging"
	"quantum/shared/protocol_v1"
	"sync"

	"github.com/golang/protobuf/proto"
)

// Handler ...
type Handler struct {
	logger     logging.Logger
	references map[string]Resource
	mutex      *sync.Mutex
}

// NewHandler ...
func NewHandler() *Handler {
	return &Handler{
		logger: &logging.DefaultLogger{},
		references: map[string]Resource{
			"core.FileSystem": &types.Core_FileSystem{},
		},
		mutex: &sync.Mutex{},
	}
}

// Handle ...
func (t *Handler) Handle(env *protocol_v1.Envelope, w io.WriteCloser) {
	defer func() {
		t.closectx(env, w)
	}()

	switch env.Type {
	case protocol_v1.Envelope_CallMember:
		t.callMember(env, w)
		break
	}
}

// callMember ...
func (t *Handler) callMember(env *protocol_v1.Envelope, w io.WriteCloser) {
	call := &protocol_v1.CallMember{}

	if err := proto.Unmarshal(env.Data, call); err != nil {
		logger.Error(err.Error())
		return
	}

	t.mutex.Lock()
	ref, ok := t.references[call.Reference]
	t.mutex.Unlock()

	if !ok {
		t.logger.Error("No such context found %v", env.Context)
		return
	}

	ref.Call(call.Name, call.Params, env, w)
}

func (t *Handler) closectx(env *protocol_v1.Envelope, w io.WriteCloser) {
	defer func() {
		w.Close()
	}()

	close := &protocol_v1.Close{}
	closeb, err := proto.Marshal(close)

	if err != nil {
		t.logger.Error(err.Error())
		return
	}

	env.Type = protocol_v1.Envelope_Close
	env.Data = closeb

	envb, err := proto.Marshal(env)

	if err != nil {
		t.logger.Error(err.Error())
		return
	}

	if _, err := w.Write(envb); err != nil {
		t.logger.Error(err.Error())
		return
	}
}
