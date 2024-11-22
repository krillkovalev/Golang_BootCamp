package server

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

const portNum string = ":8888"
const keyServerAddr = "127.0.0.1"

func Homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Api")
}

func TurnOn() {

	http.HandleFunc("/api/places", Homepage)

	err := http.ListenAndServe(portNum, nil)
	if err != nil {
		log.Fatal(err)
	}
}

func writeResponseBadRequest(w http.ResponseWriter, response interface{}) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	writeResponse()
}

func writeResponse(w http.ResponseWriter)

func getRoot(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	hasFirst := r.URL.Query().Has("first")
	first := r.URL.Query().Get("first")
	hasSecond := r.URL.Query().Has("second")
	second := r.URL.Query().Get("second")

	fmt.Printf("%s: got / request. first(%t)=%s, second(%t)=%s\n",
		ctx.Value(keyServerAddr),
		hasFirst, first,
		hasSecond, second)
}
