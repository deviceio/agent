package filesystem

import "github.com/deviceio/agent/transport"

func init() {
	fs := &filesystem{}
	transport.Router.HandleFunc("/filesystem/read", fs.read).Methods("POST")
}
