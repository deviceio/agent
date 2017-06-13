package filesystem

import (
	"github.com/deviceio/shared/logging"
	"github.com/gorilla/mux"
)

func RegisterRoutes(router *mux.Router) {
	r := &root{
		logger: &logging.DefaultLogger{},
	}

	router.HandleFunc("/filesystem", r.get).Methods("GET")
	router.HandleFunc("/filesystem/read", r.read).Methods("POST")
	router.HandleFunc("/filesystem/write", r.write).Methods("POST")
}
