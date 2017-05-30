package process

import (
	"encoding/json"
	"net/http"

	"github.com/deviceio/hmapi"
)

type process struct {
}

func (t *process) get(rw http.ResponseWriter, r *http.Request) {
	parentPath := r.Header.Get("X-Deviceio-Parent-Path")

	resource := &hmapi.Resource{
		Forms: map[string]*hmapi.Form{
			"new": &hmapi.Form{
				Action:  parentPath + "/process/new",
				Method:  "POST",
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
	}

	rw.Header().Set("Content-Type", hmapi.MediaTypeJSON.String())
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&resource)
}

func (t *process) new(rw http.ResponseWriter, r *http.Request) {

}
