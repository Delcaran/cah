// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embed"
	"flag"
	"log"
	"net/http"
	"text/template"
	"time"

	"github.com/delcaran/cah/db"
	"github.com/delcaran/cah/game"
)

//go:embed template/*
var content embed.FS

//go:embed static/*
var static embed.FS

var templates = template.Must(template.ParseFS(content, "template/*.html"))
var addr = flag.String("addr", ":8080", "http service address")

type PageContent struct {
	CurrentPlayer    *game.Player
	CurrentBlackCard *db.BlackCard
	Sets             map[string]*[]db.Set
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	game_status := game.GetGame()
	pc := PageContent{}
	r.ParseForm()
	if len(game_status.Players) > 0 {
		// game runnning
		pc.CurrentBlackCard = game_status.Black_Card
		if r.PostForm.Has("player_name") {
			// another player is joining
			game.Join(r.FormValue("player_name"))
			pc.CurrentPlayer = &game_status.Players[len(game_status.Players)-1]
			log.Printf("Player %d : %s \n", pc.CurrentPlayer.ID, pc.CurrentPlayer.Name)
		}
	} else {
		// new game
		if r.PostForm.Has("sets") {
			// first player is joining
			game.Init(r.FormValue("lang"), r.Form["sets"], r.FormValue("player_name"))
			pc.CurrentBlackCard = game_status.Black_Card
			pc.CurrentPlayer = &game_status.Players[0]
			log.Printf("CZAR %d : %s \n", pc.CurrentPlayer.ID, pc.CurrentPlayer.Name)
		} else {
			// no player yet, show sets selection
			pc.CurrentBlackCard = nil
			pc.CurrentPlayer = nil
			pc.Sets = make(map[string]*[]db.Set)
			dbs := game.Load()
			for lang, db := range *dbs {
				pc.Sets[lang] = &db.Sets
			}
		}
	}

	templates.ExecuteTemplate(w, "home.html", pc)
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.Handle("/static/", http.FileServerFS(static))
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
