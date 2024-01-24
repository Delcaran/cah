package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

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

func new_game_handler(w http.ResponseWriter, r *http.Request) {
    database, _ := db.Load("eng")
	err := templates.ExecuteTemplate(w, "new.html", database)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func card_selection_handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}
	fmt.Fprintf(w, "Post from website! r.PostFrom = %v\n", r.PostForm)
	//player_name := r.FormValue("name") // only one element
    selected_sets_str := r.Form["sets"] // array of strings
	var selected_sets []int
	for _, s := range selected_sets_str {
		v, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		selected_sets = append(selected_sets, v)
	}
	loaded_cards, err := db.GetSelectedCards(selected_sets)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
	err = templates.ExecuteTemplate(w, "test.html", loaded_cards)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
    http.HandleFunc("/", index_handler)
	http.HandleFunc("/new/", new_game_handler)
	http.HandleFunc("/select_sets/", card_selection_handler)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
