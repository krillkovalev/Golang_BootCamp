package server

import (
	"Interface/db"
	"Interface/types"
	//"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math"
	"net/http"
	"path/filepath"
	"strconv"
)

func HandlePlaces(w http.ResponseWriter, r *http.Request) {

	s := r.URL.Query().Get("page")
	page, err := strconv.Atoi(s)
	if err != nil {
		log.Fatalf("Error to parse query: %s", err)
	}
	if page < 1 {
		w.WriteHeader(http.StatusBadRequest)
	}

	limit := 10
	offset := (page - 1) * limit
	places, total, err := db.GetPlaces(limit, offset)
	fmt.Println(places)

	totalpages := int(math.Round(float64(total) / float64(limit)))

	data := struct {
		Total      int
		Places     []types.Place
		TotalPages int
		Page       int
	}{
		Total:      total,
		Places:     places,
		TotalPages: totalpages,
		Page:       page,
	}
	path := filepath.Join("server", "template.html")
	tmpl, err := template.New(filepath.Base(path)).Funcs(template.FuncMap{
		"add":   func(a, b int) int { return a + b },
		"minus": func(a, b int) int { return a - b },
	}).ParseFiles(path)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Fatal(err)
	}

}
