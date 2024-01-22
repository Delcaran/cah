package main

import (
	_ "embed"

	"fmt"
	"log"
	"net/http"

	"github.com/delcaran/cah/db"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func db_load(w http.ResponseWriter, r *http.Request) {
    database, _ := db.Load("eng")
	var content string
	for set_num, set := range database.Sets {
		content = fmt.Sprintf("%s<li>[%d] %s</li>\n", content, set_num, set.Name)
	}
    fmt.Fprintf(w, "", content)
}

func main() {
    http.HandleFunc("/", handler)
	http.HandleFunc("/load", db_load)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
