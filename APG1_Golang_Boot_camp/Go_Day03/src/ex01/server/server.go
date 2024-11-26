package server

import (
	//"fmt"
	//"html/template"
	"Interface/db"
	"log"
	"net/http"
	//"net/url"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8"
)

func HandlePlaces(w http.ResponseWriter, r *http.Request) {
	client, err := elasticsearch.NewDefaultClient()

	client.Count(
		client.Count.WithIndex("places"),
	)

	s := r.URL.Query().Get("page")
	page, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal(err)
	}
	if page < 1 {
		w.WriteHeader(http.StatusBadRequest)
	}

	limit := 10
	offset := (page - 1) * 10
	places, cnt, err := db.GetPlaces(limit, offset)

	//totalpages :=

}
