package main

import (
	"Interface/db"
	"encoding/json"
	"log"
	"net.http"
	"net/http"
	"strconv"
)

type Server struct {
	client *Client
}

func main() {
	client := NewClient("http://localhost:8888")
	if err := client.CheckHealth(); err != nil {
		log.Fatalf("failed to health check", err)
	}
	if err := client.CreateIndex(); err != nil {
		log.Fatal("failed create index: ", err)
	}

	server := Server{client: client}
	http.HandleFunc("/insert", server.InsertDataHandler)
	http.HandleFunc("/update", server.InsertDataHandler)
	http.HandleFunc("/delete", server.InsertDataHandler)
	http.HandleFunc("/search", server.InsertDataHandler)
	http.HandleFunc("/health", server.InsertDataHandler)

	http.ListenAndServe(":8080", nil)
}
