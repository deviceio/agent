package filesystem

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
)

type filesystem struct {
}

func (t *filesystem) read(w http.ResponseWriter, r *http.Request) {
	type model struct {
		Path string `json:"path"`
	}

	var m *model

	err := json.NewDecoder(r.Body).Decode(&m)

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	file, err := os.Open(m.Path)
	defer file.Close()

	if err != nil {
		w.WriteHeader(400)
		w.Write([]byte(err.Error()))
		return
	}

	io.Copy(w, file)
}
