package router

import (
	"restore-photos/controller"

	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/generate", controller.Generate).Methods("POST")

	return router
}
