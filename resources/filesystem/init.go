package filesystem

import (
	"sync"

	"github.com/deviceio/agent/transport"
	"github.com/deviceio/shared/logging"
)

func init() {
	fs := &filesystem{
		logger: &logging.DefaultLogger{},
		handles: &handleCollection{
			items: []*handle{},
			mutex: &sync.Mutex{},
		},
	}

	transport.Router.HandleFunc("/filesystem", fs.get).Methods("GET")
	transport.Router.HandleFunc("/filesystem/handle", fs.handles.get).Methods("GET")
	transport.Router.HandleFunc("/filesystem/handle/{index}", fs.handles.getitem).Methods("GET")
	transport.Router.HandleFunc("/filesystem/open", fs.open).Methods("POST")
	transport.Router.HandleFunc("/filesystem/read", fs.read).Methods("POST")
	transport.Router.HandleFunc("/filesystem/write", fs.write).Methods("POST")
}
