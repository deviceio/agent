package filesystem

import (
	"github.com/deviceio/agent/transport"
	"github.com/deviceio/shared/logging"
)

func init() {
	fs := &filesystem{
		logger: &logging.DefaultLogger{},
	}

	transport.Router.HandleFunc("/filesystem", fs.get).Methods("GET")
	transport.Router.HandleFunc("/filesystem/read", fs.read).Methods("POST")
	transport.Router.HandleFunc("/filesystem/write", fs.write).Methods("POST")
}
