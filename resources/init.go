package resources

import "github.com/deviceio/agent/transport"

func init() {
	r := &root{}
	transport.Router.HandleFunc("/", r.get).Methods("GET")
}
