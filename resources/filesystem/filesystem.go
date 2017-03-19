package filesystem

import (
	"io"
	"net/http"
	"os"
	"strconv"

	"github.com/deviceio/shared/logging"
)

type filesystem struct {
	logger logging.Logger
}

func (t *filesystem) read(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Trailer", "Error")
	w.WriteHeader(200)

	var file *os.File
	var err error
	var count int64 = -1
	var offset int64
	var offsetAt int
	var path string

	args := map[string]string{
		"count":    r.Header.Get("X-Count"),
		"offset":   r.Header.Get("X-Offset"),
		"offsetAt": r.Header.Get("X-OffsetAt"),
		"path":     r.Header.Get("X-Path"),
	}

	if count, err = strconv.ParseInt(args["count"], 10, 64); args["count"] != "" && err != nil {
		w.Write([]byte(" "))
		w.Header().Set("Error", err.Error())
		return
	}

	if offset, err = strconv.ParseInt(args["offset"], 10, 64); args["offset"] != "" && err != nil {
		w.Write([]byte(" "))
		w.Header().Set("Error", err.Error())
		return
	}

	if offsetAt, err = strconv.Atoi(args["offsetAt"]); args["offsetAt"] != "" && err != nil {
		w.Write([]byte(" "))
		w.Header().Set("Error", err.Error())
		return
	}

	path = args["path"]

	if file, err = os.Open(path); err != nil {
		w.Write([]byte(" "))
		w.Header().Set("Error", err.Error())
		return
	}
	defer file.Close()

	if offsetAt == 0 {
		if _, err = file.Seek(offset, offsetAt); err != nil {
			w.Write([]byte(" "))
			w.Header().Set("Error", err.Error())
			return
		}
	} else if offsetAt == 1 {
		if _, err = file.Seek(offset, 2); err != nil {
			w.Write([]byte(" "))
			w.Header().Set("Error", err.Error())
			return
		}
	} else {
		w.Write([]byte(" "))
		w.Header().Set("Error", "Unknown offsetAt code")
		return
	}

	var curcnt int64
	buf := make([]byte, 250000)

	if n, err := io.CopyBuffer(w, file, buf); err != nil {
		curcnt += n

		if err != io.EOF {
			w.Write([]byte(" "))
			w.Header().Set("Error", err.Error())
			return
		}

		if count > 0 && curcnt > count {
			return
		}
	}
}
