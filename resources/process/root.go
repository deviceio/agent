package process

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"sync"

	"github.com/deviceio/hmapi"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type root struct {
	itemsmu *sync.Mutex
	items   map[string]*process
}

func (t *root) get(rw http.ResponseWriter, r *http.Request) {
	parentPath := r.Header.Get("X-Deviceio-Parent-Path")

	resource := &hmapi.Resource{
		Forms: map[string]*hmapi.Form{
			"create-process": &hmapi.Form{
				Action:  parentPath + "/process",
				Method:  hmapi.POST,
				Enctype: hmapi.MediaTypeMultipartFormData,
				Fields: []*hmapi.FormField{
					&hmapi.FormField{
						Name:     "cmd",
						Type:     hmapi.MediaTypeHMAPIString,
						Required: true,
					},
					&hmapi.FormField{
						Name:     "args",
						Type:     hmapi.MediaTypeHMAPIString,
						Multiple: true,
					},
				},
			},
		},
		Content: map[string]*hmapi.Content{},
		Links:   map[string]*hmapi.Link{},
	}

	t.itemsmu.Lock()
	defer t.itemsmu.Unlock()

	resource.Content["process-list"] = &hmapi.Content{
		Type:  hmapi.MediaTypeJSON,
		Value: t.items,
	}

	for id := range t.items {
		resource.Links[id] = &hmapi.Link{
			Href: parentPath + "/process/" + id,
			Type: hmapi.MediaTypeHMAPIResource,
		}
	}

	rw.Header().Set("Content-Type", hmapi.MediaTypeJSON.String())
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&resource)
}

func (t *root) createProcess(rw http.ResponseWriter, r *http.Request) {
	parentPath := r.Header.Get("X-Deviceio-Parent-Path")
	formReader, err := r.MultipartReader()

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	form, err := formReader.ReadForm(250000)

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	cmds, ok := form.Value["cmd"]

	if !ok {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("cmd not supplied"))
		return
	}

	if len(cmds) != 1 {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("only one cmd field should be supplied"))
		return
	}

	cmdstr := cmds[0]

	args, ok := form.Value["arg"]

	if !ok {
		args = []string{}
	}

	t.itemsmu.Lock()
	defer t.itemsmu.Unlock()

	procid, _ := uuid.NewRandom()

	ctx, cancel := context.WithCancel(context.Background())

	proc := &process{
		id:     procid.String(),
		cmd:    exec.CommandContext(ctx, cmdstr, args...),
		mu:     &sync.Mutex{},
		ctx:    ctx,
		cancel: cancel,
	}

	t.items[procid.String()] = proc

	rw.Header().Set("Location", fmt.Sprintf(
		"%v/%v/%v",
		parentPath,
		"process",
		procid.String(),
	))

	rw.WriteHeader(http.StatusCreated)
	rw.Write([]byte(""))
}

func (t *root) getProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc-id"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.get(rw, r)
}

func (t *root) deleteProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc-id"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.stop(rw, r)
	delete(t.items, vars["proc-id"])

	rw.WriteHeader(http.StatusNoContent)
}

func (t *root) startProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc-id"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.start(rw, r)
}

func (t *root) stopProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc-id"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.stop(rw, r)
}

func (t *root) stdinProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc-id"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.stdin(rw, r)
}

func (t *root) stdoutProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc-id"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.stdout(rw, r)
}

func (t *root) stderrProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc-id"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.stderr(rw, r)
}

func (t *root) findProcessItem(processid string) *process {
	t.itemsmu.Lock()
	defer t.itemsmu.Unlock()

	proc, ok := t.items[processid]

	if !ok {
		return nil
	}

	return proc
}
