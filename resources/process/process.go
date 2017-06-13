package process

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"strings"
	"sync"

	"mime/multipart"

	"github.com/Sirupsen/logrus"
	"github.com/deviceio/hmapi"
	"github.com/google/uuid"
	"github.com/palantir/stacktrace"
)

type processReaders struct {
	items []*io.PipeWriter
	*sync.RWMutex
}

type process struct {
	id            string
	cmd           *exec.Cmd
	mu            *sync.Mutex
	ctx           context.Context
	cancel        context.CancelFunc
	started       bool
	stdoutPipe    io.ReadCloser
	stderrPipe    io.ReadCloser
	stdinPipe     io.WriteCloser
	stdoutReaders *processReaders
	stderrReaders *processReaders
}

func newProcess(cmd string, args []string) (*process, error) {
	procid, err := uuid.NewRandom()

	if err != nil {
		return nil, stacktrace.Propagate(err, "error generating proc id")
	}

	ctx, cancel := context.WithCancel(context.Background())

	proc := &process{
		id:     strings.ToLower(procid.String()),
		cmd:    exec.CommandContext(ctx, cmd, args...),
		mu:     &sync.Mutex{},
		ctx:    ctx,
		cancel: cancel,
		stdoutReaders: &processReaders{
			RWMutex: &sync.RWMutex{},
			items:   []*io.PipeWriter{},
		},
		stderrReaders: &processReaders{
			RWMutex: &sync.RWMutex{},
			items:   []*io.PipeWriter{},
		},
	}

	stdoutPipe, err := proc.cmd.StdoutPipe()

	if err != nil {
		return nil, stacktrace.Propagate(err, "error creating stdout pipe")
	}

	stderrPipe, err := proc.cmd.StderrPipe()

	if err != nil {
		return nil, stacktrace.Propagate(err, "error creating stderr pipe")
	}

	stdinPipe, err := proc.cmd.StdinPipe()

	if err != nil {
		return nil, stacktrace.Propagate(err, "error creating stdin pipe")
	}

	proc.stdoutPipe = stdoutPipe
	proc.stderrPipe = stderrPipe
	proc.stdinPipe = stdinPipe

	return proc, nil
}

func (t *process) get(rw http.ResponseWriter, r *http.Request) {
	parentPath := r.Header.Get("X-Deviceio-Parent-Path")

	resource := &hmapi.Resource{
		Links: map[string]*hmapi.Link{
			"self": &hmapi.Link{
				Href: parentPath + "/process/" + t.id,
			},
			"parent": &hmapi.Link{
				Href: parentPath + "/process",
			},
			"stdout": &hmapi.Link{
				Href: fmt.Sprintf("%v/process/%v/stdout", parentPath, t.id),
				Type: hmapi.MediaTypeOctetStream,
			},
			"stderr": &hmapi.Link{
				Href: fmt.Sprintf("%v/process/%v/stderr", parentPath, t.id),
				Type: hmapi.MediaTypeOctetStream,
			},
		},
		Forms: map[string]*hmapi.Form{
			"delete": &hmapi.Form{
				Action:  fmt.Sprintf("%v/process/%v", parentPath, t.id),
				Enctype: hmapi.MediaTypeMultipartFormData,
				Method:  hmapi.DELETE,
			},
			"stdin": &hmapi.Form{
				Action:  fmt.Sprintf("%v/process/%v/stdin", parentPath, t.id),
				Method:  hmapi.POST,
				Enctype: hmapi.MediaTypeMultipartFormData,
				Fields: []*hmapi.FormField{
					&hmapi.FormField{
						Name:     "data",
						Type:     hmapi.MediaTypeOctetStream,
						Required: true,
					},
				},
			},
		},
		Content: map[string]*hmapi.Content{
			"id": &hmapi.Content{
				Type:  hmapi.MediaTypeHMAPIString,
				Value: t.id,
			},
			"cmd": &hmapi.Content{
				Type:  hmapi.MediaTypeHMAPIString,
				Value: t.cmd.Path,
			},
			"args": &hmapi.Content{
				Type:  hmapi.MediaTypeJSON,
				Value: t.cmd.Args,
			},
			"started": &hmapi.Content{
				Type:  hmapi.MediaTypeHMAPIBoolean,
				Value: t.started,
			},
		},
	}

	if t.started {
		resource.Forms["stop"] = &hmapi.Form{
			Action:  fmt.Sprintf("%v/process/%v/stop", parentPath, t.id),
			Enctype: hmapi.MediaTypeMultipartFormData,
			Method:  hmapi.POST,
		}
	} else {
		resource.Forms["start"] = &hmapi.Form{
			Action:  fmt.Sprintf("%v/process/%v/start", parentPath, t.id),
			Enctype: hmapi.MediaTypeMultipartFormData,
			Method:  hmapi.POST,
		}
	}

	rw.Header().Set("Content-Type", hmapi.MediaTypeJSON.String())
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&resource)
}

func (t *process) start(rw http.ResponseWriter, r *http.Request) {
	if t.started {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("process already started"))
	}

	err := t.cmd.Start()

	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}

	stdoutbuf := make([]byte, 250000)
	stderrbuf := make([]byte, 250000)

	go func(t *process) {
		for {
			stdoutn, err := t.stdoutPipe.Read(stdoutbuf)

			if stdoutn > 0 {
				data := stdoutbuf[:stdoutn]

				t.stdoutReaders.RLock()
				for _, reader := range t.stdoutReaders.items {
					if _, err := reader.Write(data); err != nil {
						reader.Close()
					}
				}
				t.stdoutReaders.RUnlock()
			}

			if err != nil {
				t.stdoutReaders.RLock()
				for _, reader := range t.stdoutReaders.items {
					reader.Close()
				}
				t.stdoutReaders.RUnlock()
				break
			}
		}

		err := t.cmd.Wait()

		if err != nil {
			logrus.WithField("error", err.Error()).Error("cmd exited with error")
		}

		t.cancel()
	}(t)

	go func(t *process) {
		for {
			stderrn, err := t.stderrPipe.Read(stderrbuf)

			if stderrn > 0 {
				data := stderrbuf[:stderrn]

				t.stderrReaders.RLock()
				for _, reader := range t.stderrReaders.items {
					if _, err := reader.Write(data); err != nil {
						reader.Close()
					}
				}
				t.stderrReaders.RUnlock()
			}

			if err != nil {
				t.stderrReaders.RLock()
				for _, reader := range t.stderrReaders.items {
					reader.Close()
				}
				t.stderrReaders.RUnlock()
				break
			}
		}
	}(t)

	t.started = true
}

func (t *process) stop(rw http.ResponseWriter, r *http.Request) {
	if !t.started {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("process not started"))
	}

	t.cancel()
}

func (t *process) stdin(rw http.ResponseWriter, r *http.Request) {
	ps := t.cmd.ProcessState

	if ps != nil && ps.Exited() {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("cannot supply stdin process has already exited"))
		return
	}

	form := multipart.NewReader(r.Body, hmapi.MultipartFormDataBoundry)
	data, err := form.NextPart()

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	if data.FormName() != "data" {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("field 'data' not supplied"))
		return
	}

	close := rw.(http.CloseNotifier).CloseNotify()

	buf := make([]byte, 250000)
	chdata := make(chan []byte)
	cherr := make(chan error)
	done := make(chan bool)

	go func() {
		for {
			n, err := data.Read(buf)

			if n > 0 {
				chdata <- buf[:n]
			}

			if err != nil && err != io.EOF {
				cherr <- err
			}

			if err == io.EOF {
				done <- true
				return
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case <-t.ctx.Done():
			return
		case <-close:
			return
		case data := <-chdata:
			t.stdinPipe.Write(data)
		case err := <-cherr:
			rw.WriteHeader(http.StatusInternalServerError)
			rw.Write([]byte(err.Error()))
			return
		}
	}
}

func (t *process) stdout(rw http.ResponseWriter, r *http.Request) {
	ps := t.cmd.ProcessState

	if ps != nil && ps.Exited() {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("cannot supply stdout process has already exited"))
		return
	}

	flush := rw.(http.Flusher).Flush

	t.stdoutReaders.Lock()
	piper, pipew := io.Pipe()
	t.stdoutReaders.items = append(t.stdoutReaders.items, pipew)
	t.stdoutReaders.Unlock()

	rw.Header().Set("Trailer", "Error")
	rw.Header().Set("Content-Type", hmapi.MediaTypeOctetStream.String())
	rw.WriteHeader(http.StatusOK)
	flush()

	buf := make([]byte, 250000)

	for {
		n, err := piper.Read(buf)

		if n > 0 {
			rw.Write(buf[:n])
			flush()
		}

		if err != nil {
			if err != io.EOF {
				rw.Header().Set("Error", err.Error())
			}
			return
		}
	}
}

func (t *process) stderr(rw http.ResponseWriter, r *http.Request) {
	ps := t.cmd.ProcessState

	if ps != nil && ps.Exited() {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte("cannot supply stderr process has already exited"))
		return
	}

	flush := rw.(http.Flusher).Flush

	t.stderrReaders.Lock()
	piper, pipew := io.Pipe()
	t.stderrReaders.items = append(t.stderrReaders.items, pipew)
	t.stderrReaders.Unlock()

	rw.Header().Set("Trailer", "Error")
	rw.Header().Set("Content-Type", hmapi.MediaTypeOctetStream.String())
	rw.WriteHeader(http.StatusOK)
	flush()

	buf := make([]byte, 250000)

	for {
		n, err := piper.Read(buf)

		if n > 0 {
			rw.Write(buf[:n])
			flush()
		}

		if err != nil {
			if err != io.EOF {
				rw.Header().Set("Error", err.Error())
			}
			return
		}
	}
}
