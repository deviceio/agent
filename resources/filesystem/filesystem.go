package filesystem

import (
	"io"
	"net/http"
	"os"

	"github.com/deviceio/shared/logging"
)

type filesystem struct {
	logger logging.Logger
}

func (t *filesystem) read(w http.ResponseWriter, r *http.Request) {
	//flusher, _ := w.(http.Flusher)

	path := r.Header.Get("X-Path")

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	buf := make([]byte, 250000)

	if _, err := io.CopyBuffer(w, file, buf); err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
	}
}
