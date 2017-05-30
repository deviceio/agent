package process

import "github.com/deviceio/agent/transport"

func init() {
	proc := &process{}

	transport.Router.HandleFunc("/process", proc.get).Methods("GET")
	transport.Router.HandleFunc("/process/new", proc.new).Methods("POST")
}
