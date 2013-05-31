package main

import (
	"github.com/ManuelKiessling/infmgmt-backend/interfaces"
	"github.com/gorilla/mux"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/machines", interfaces.MachinesRequesthandler)
	http.Handle("/", r)
	http.ListenAndServe(":8080", nil)
}
