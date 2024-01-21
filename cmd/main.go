package main

import (
	"net/http"

	"server-sent-events-example/pkg"

	"github.com/gorilla/mux"
)

func main() {
	addrKeeper := pkg.NewAddrKeeper()
	router := mux.NewRouter()
	router.Handle("/users/{uid}/messages", pkg.PullMsgFactory(addrKeeper)).Methods("GET")
	router.Handle("/users/{uid}/messages", pkg.SendMsgFactory(addrKeeper)).Methods("POST")
	http.ListenAndServe(":8080", router)
}
