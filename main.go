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
var templates = template.Must(template.ParseFS(content, "template/*.html"))

var addr = flag.String("addr", ":8080", "http service address")

type PageContent struct {
	CurrentPlayer    *game.Player
	CurrentBlackCard *db.BlackCard
}

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	game_status := game.GetGame()
	pc := PageContent{}
	if len(game_status.Players) > 0 {
		// game runnning
		//TODO fill game info (player ID or name where?)
	} else {
		// new game
		r.ParseForm()
		if r.PostForm.Has("sets") {
			// match can begin
			game.Init(r.Form["sets"], r.FormValue("player_name"))
			pc.CurrentBlackCard = game_status.Black_Card
			pc.CurrentPlayer = &game_status.Players[0]
		} else if r.PostForm.Has("player_name") {
			// a player is joining
			game.Join(r.FormValue("player_name"))
			pc.CurrentBlackCard = game_status.Black_Card
			pc.CurrentPlayer = &game_status.Players[len(game_status.Players)-1]
		} else {
			// sets selection and first player name
			pc.CurrentBlackCard = nil
			pc.CurrentPlayer = nil
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
	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
