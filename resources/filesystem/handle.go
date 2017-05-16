package filesystem

import (
	"encoding/json"
	"net/http"
	"os"
	"sync"

	"github.com/deviceio/hmapi"
)

type handle struct {
	id   string
	file *os.File
}

type handleCollection struct {
	items []*handle
	mutex *sync.Mutex
}

func (t *handleCollection) get(rw http.ResponseWriter, r *http.Request) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	parentPath := r.Header.Get("X-Deviceio-Parent-Path")

	resource := &hmapi.Resource{
		Links: map[string]*hmapi.Link{
			"parent": &hmapi.Link{
				Href: parentPath + "/filesystem",
				Type: hmapi.MediaTypeJSON,
			},
		},
		Content: map[string]*hmapi.Content{
			"count": &hmapi.Content{
				Type:  hmapi.MediaTypeHMAPIInt,
				Value: len(t.items),
			},
		},
	}

	for _, handle := range t.items {
		resource.Links[handle.id] = &hmapi.Link{
			Href: parentPath + "/filesystem/handle/" + handle.id,
			Type: hmapi.MediaTypeJSON,
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&resource)
}

func (t *handleCollection) getitem(rw http.ResponseWriter, r *http.Request) {

}
