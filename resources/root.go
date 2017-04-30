package resources

import (
	"encoding/json"
	"net/http"
	"os"
	"runtime"

	"github.com/deviceio/hmapi"
	"github.com/spf13/viper"
)

type root struct {
}

func (t *root) get(rw http.ResponseWriter, r *http.Request) {
	resource := &hmapi.Resource{
		Links: map[string]*hmapi.Link{
			"system": &hmapi.Link{
				Type: hmapi.MediaTypeJSON,
				Href: "/system",
			},
			"filesystem": &hmapi.Link{
				Type: hmapi.MediaTypeJSON,
				Href: "/filesystem",
			},
		},
		Content: map[string]*hmapi.Content{
			"id": &hmapi.Content{
				Type:  hmapi.MediaTypeHMAPIString,
				Value: viper.GetString("id"),
			},
			"hostname": &hmapi.Content{
				Type: hmapi.MediaTypeHMAPIString,
				Value: (func() string {
					hostname, _ := os.Hostname()
					return hostname
				})(),
			},
			"architecture": &hmapi.Content{
				Type:  hmapi.MediaTypeHMAPIString,
				Value: runtime.GOARCH,
			},
			"platform": &hmapi.Content{
				Type:  hmapi.MediaTypeHMAPIString,
				Value: runtime.GOOS,
			},
		},
	}

	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(http.StatusOK)
	json.NewEncoder(rw).Encode(&resource)
}
