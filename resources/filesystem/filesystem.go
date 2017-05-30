package filesystem

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"

	"io/ioutil"

	"github.com/deviceio/hmapi"
	"github.com/deviceio/shared/logging"
)

type filesystem struct {
	logger logging.Logger
}

func (t *filesystem) get(rw http.ResponseWriter, r *http.Request) {
	parentPath := r.Header.Get("X-Deviceio-Parent-Path")

	resource := &hmapi.Resource{
		Links: map[string]*hmapi.Link{},
		Forms: map[string]*hmapi.Form{
			"read": &hmapi.Form{
				Action:  parentPath + "/filesystem/read",
				Method:  "POST",
				Type:    hmapi.MediaTypeOctetStream,
				Enctype: hmapi.MediaTypeMultipartFormData,
				Fields: []*hmapi.FormField{
					&hmapi.FormField{
						Name:     "path",
						Type:     hmapi.MediaTypeHMAPIString,
						Required: true,
					},
					&hmapi.FormField{
						Name:     "offset",
						Type:     hmapi.MediaTypeHMAPIInt,
						Required: false,
						Value:    0,
					},
					&hmapi.FormField{
						Name:     "offsetAt",
						Type:     hmapi.MediaTypeHMAPIInt,
						Required: false,
						Value:    0,
					},
					&hmapi.FormField{
						Name:     "count",
						Type:     hmapi.MediaTypeHMAPIInt,
						Required: false,
						Value:    -1,
					},
				},
			},
			"write": &hmapi.Form{
				Action:  parentPath + "/filesystem/write",
				Method:  "POST",
				Type:    hmapi.MediaTypeHMAPIInt,
				Enctype: hmapi.MediaTypeMultipartFormData,
				Fields: []*hmapi.FormField{
					&hmapi.FormField{
						Name:     "path",
						Type:     hmapi.MediaTypeHMAPIString,
						Required: true,
					},
					&hmapi.FormField{
						Name:     "append",
						Type:     hmapi.MediaTypeHMAPIBoolean,
						Required: false,
						Value:    false,
					},
					&hmapi.FormField{
						Name:     "data",
						Type:     hmapi.MediaTypeOctetStream,
						Required: true,
					},
				},
			},
		},
		Content: map[string]*hmapi.Content{},
	}

	rw.Header().Set("Content-Type", hmapi.MediaTypeJSON.String())
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&resource)
}

func (t *filesystem) read(w http.ResponseWriter, r *http.Request) {
	var file *os.File
	var err error
	var count int64 = -1
	var offset int64
	var offsetAt int
	var path string

	err = r.ParseMultipartForm(4096)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	args := map[string]string{
		"count":    r.FormValue("count"),
		"offset":   r.FormValue("offset"),
		"offsetAt": r.FormValue("offsetAt"),
		"path":     r.FormValue("path"),
	}

	if args["path"] == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Path must be supplied"))
		return
	}

	if count, err = strconv.ParseInt(args["count"], 10, 64); args["count"] != "" && err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if offset, err = strconv.ParseInt(args["offset"], 10, 64); args["offset"] != "" && err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	if offsetAt, err = strconv.Atoi(args["offsetAt"]); args["offsetAt"] != "" && err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	path = args["path"]

	if file, err = os.Open(path); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}
	defer file.Close()

	if offsetAt == 0 {
		if _, err = file.Seek(offset, offsetAt); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	} else if offsetAt == 1 {
		if _, err = file.Seek(offset, 2); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}
	} else {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unknown OffsetAt Code"))
		return
	}

	var currcnt int64
	buf := make([]byte, 250000)

	// start chunk writing
	w.Header().Set("Trailer", "Error")
	w.Header().Set("Content-Type", hmapi.MediaTypeOctetStream.String())
	w.WriteHeader(200)

	flusher, flusherok := w.(http.Flusher)

	for {
		nr, rerr := file.Read(buf)

		if rerr != nil && rerr != io.EOF {
			w.Write([]byte(""))
			w.Header().Set("Error", err.Error())
			return
		}

		if rerr != nil && rerr == io.EOF {
			w.Write([]byte(""))
			w.Header().Set("Error", "")
			return
		}

		if nr == 0 {
			continue
		}

		nw, werr := w.Write(buf[:nr])

		if flusherok {
			flusher.Flush()
		}

		if werr != nil {
			w.Write([]byte(""))
			w.Header().Set("Error", err.Error())
			return
		}

		currcnt += int64(nw)

		if count > 0 && currcnt >= count {
			w.Write([]byte(""))
			w.Header().Set("Error", "")
			return
		}
	}
}

func (t *filesystem) write(rw http.ResponseWriter, r *http.Request) {
	var file *os.File
	var form *multipart.Reader
	var err error
	var path string
	var appendFile = false

	if form, err = r.MultipartReader(); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		rw.Write([]byte(err.Error()))
		return
	}

	for {
		part, err := form.NextPart()

		if err == io.EOF {
			break
		} else if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte(err.Error()))
			return
		}

		if part.FormName() == "path" {
			pathb, _ := ioutil.ReadAll(part)
			path = string(pathb)
		} else if part.FormName() == "append" {
			boolb, err := ioutil.ReadAll(part)

			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(err.Error()))
				return
			}

			b, err := strconv.ParseBool(string(boolb))

			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(err.Error()))
				return
			}

			appendFile = b
		} else if part.FormName() == "data" {
			flag := os.O_RDWR | os.O_TRUNC | os.O_CREATE

			if appendFile {
				flag = os.O_RDWR | os.O_APPEND | os.O_CREATE
			}

			file, err = os.OpenFile(path, flag, 0777)
			defer file.Close()

			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(err.Error()))
				return
			}

			if _, err = io.Copy(file, part); err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				rw.Write([]byte(err.Error()))
				return
			}
		} else {
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write([]byte("Fields must be supplied in the order described by the filesystem resource"))
			return
		}
	}
}
