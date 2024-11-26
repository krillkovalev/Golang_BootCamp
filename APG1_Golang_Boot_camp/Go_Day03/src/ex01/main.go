package main

import (
	//"Interface/db"
	"Interface/server"
	"net/http"
)

func main() {

	http.HandleFunc("/", server.HandlePlaces)

	http.ListenAndServe(":8888", nil)
}
