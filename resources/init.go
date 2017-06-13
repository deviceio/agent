package resources

import "github.com/gorilla/mux"

func RegisterRoutes(router *mux.Router) {
	r := &root{}
	router.HandleFunc("/", r.get).Methods("GET")
	router.HandleFunc("/chunk-response-test", r.chunkedResponseTest).Methods("GET")
}
