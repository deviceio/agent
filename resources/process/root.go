package process

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/Sirupsen/logrus"
	"github.com/deviceio/hmapi"
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
			"create": &hmapi.Form{
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
						Name:     "arg",
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

	proc, err := newProcess(cmdstr, args)

	if err != nil {
		logrus.WithField("error", err.Error()).Error("failed to create new process")
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte("failed to create process. check agent logs for more details"))
		return
	}

	t.items[proc.id] = proc

	rw.Header().Set("Location", fmt.Sprintf(
		"%v/%v/%v",
		parentPath,
		"process",
		proc.id,
	))

	rw.WriteHeader(http.StatusCreated)
}

func (t *root) getProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.get(rw, r)
}

func (t *root) deleteProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.cancel()

	delete(t.items, vars["proc"])

	rw.WriteHeader(http.StatusOK)
}

func (t *root) startProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.start(rw, r)
}

func (t *root) stopProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.stop(rw, r)
}

func (t *root) stdinProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.stdin(rw, r)
}

func (t *root) stdoutProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc"])

	if proc == nil {
		rw.WriteHeader(http.StatusNotFound)
		rw.Write([]byte(""))
		return
	}

	proc.stdout(rw, r)
}

func (t *root) stderrProcess(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	proc := t.findProcessItem(vars["proc"])

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

	proc, ok := t.items[strings.ToLower(processid)]

	if !ok {
		return nil
	}

	return proc
}
