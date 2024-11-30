package main

import (
	"Interface/server"
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/", server.HandlePlaces)
	fmt.Println("server is running")
	http.ListenAndServe(":8888", nil)
}
