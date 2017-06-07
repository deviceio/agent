package process

import (
	"sync"

	"github.com/deviceio/agent/transport"
)

func init() {
	root := &root{
		itemsmu: &sync.Mutex{},
		items:   map[string]*process{},
	}

	transport.Router.HandleFunc("/process", root.get).Methods("GET")
	transport.Router.HandleFunc("/process", root.createProcess).Methods("POST")
	transport.Router.HandleFunc("/process/{proc-id}", root.getProcess).Methods("GET")
	transport.Router.HandleFunc("/process/{proc-id}/start", root.startProcess).Methods("POST")
	transport.Router.HandleFunc("/process/{proc-id}/stop", root.stopProcess).Methods("POST")
	transport.Router.HandleFunc("/process/{proc-id}/stdin", root.stdinProcess).Methods("POST")
	transport.Router.HandleFunc("/process/{proc-id}/stdout", root.stdoutProcess).Methods("GET")
	transport.Router.HandleFunc("/process/{proc-id}/stderr", root.stderrProcess).Methods("GET")
	transport.Router.HandleFunc("/process/{proc-id}", root.deleteProcess).Methods("DELETE")
}
