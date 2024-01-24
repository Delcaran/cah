package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/delcaran/cah/db"
)

//go:embed template/*
var content embed.FS

var templates = template.Must(template.ParseFS(content, "template/*.html"))

func index_handler(w http.ResponseWriter, r *http.Request) {
    err := templates.ExecuteTemplate(w, "main.html", nil)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func db_load(w http.ResponseWriter, r *http.Request) {
    database, _ := db.Load("eng")
	err := templates.ExecuteTemplate(w, "new.html", database)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    http.HandleFunc("/", index_handler)
	http.HandleFunc("/new", db_load)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
