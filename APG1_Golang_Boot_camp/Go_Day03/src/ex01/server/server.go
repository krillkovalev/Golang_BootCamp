package server

import (
	"Interface/db"
	"Interface/types"
	"html/template"
	"net/http"
	"path/filepath"
	"strconv"
	"fmt"
)

func HandlePlaces(w http.ResponseWriter, r *http.Request) {

	s := r.URL.Query().Get("page")
	page, err := strconv.Atoi(s)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "<html><h1>400 Bad Request</h1><p>Invalid 'page' value</p></html>")
		return
	}
	if page < 1 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "<html><h1>400 Bad Request</h1><p>Invalid 'page' value</p></html>")
		return
	}

	limit := 10
	offset := (page - 1) * limit
	places, total, err := db.GetPlaces(limit, offset)
	
	if err != nil || page > total / limit {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "<html><h1>400 Bad Request</h1><p>Invalid 'page' value</p></html>")
		return
	}

	totalpages := (total + limit - 1) / limit
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
		fmt.Fprintf(w, "<html><h1>400 Bad Request</h1><p>Invalid 'page' value</p></html>")
		return
	}
	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "<html><h1>400 Bad Request</h1><p>Invalid 'page' value</p></html>")
		return
	}

}
