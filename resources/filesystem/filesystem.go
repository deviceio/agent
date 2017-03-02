package filesystem

import (
	"io"
	"net/http"
	"os"
)

type filesystem struct {
}

func (t *filesystem) read(w http.ResponseWriter, r *http.Request) {
	flusher, _ := w.(http.Flusher)

	path := r.Header.Get("X-Path")

	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	stat, err := file.Stat()

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(200)
	w.Header().Set("Content-Length", string(stat.Size()))
	w.Header().Set("Transfer-Encoding", "chunked")
	flusher.Flush()

	buf := make([]byte, 250000)
	for {
		i, err := file.Read(buf)

		if err == io.EOF {
			w.Write([]byte(""))
			break
		}

		w.Write(buf[:i])
		flusher.Flush()
	}
}
