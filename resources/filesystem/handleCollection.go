package filesystem

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/deviceio/hmapi"
)

type handleCollection struct {
	items []*handle
	mutex *sync.Mutex
}

func (t *handleCollection) get(rw http.ResponseWriter, r *http.Request) {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	resource := &hmapi.Resource{
		Links: map[string]*hmapi.Link{
			"parent": &hmapi.Link{
				Href: "/filesystem",
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

	if len(t.items) > 0 {
		resource.Links["first"] = &hmapi.Link{
			Href: "/filesystem/handles/" + strconv.Itoa(0),
			Type: hmapi.MediaTypeJSON,
		}

		resource.Links["last"] = &hmapi.Link{
			Href: "/filesystem/handles/" + strconv.Itoa(len(t.items)-1),
			Type: hmapi.MediaTypeJSON,
		}
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&resource)
}

func (t *handleCollection) getitem(rw http.ResponseWriter, r *http.Request) {

}
