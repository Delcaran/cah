// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"embed"
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"text/template"

	"github.com/delcaran/cah/db"
	"github.com/delcaran/cah/game"
)

//go:embed template/*
var content embed.FS

//go:embed scripts/*.js
var scripts embed.FS

//go:embed styles/*.css
var styles embed.FS

var templates = template.Must(template.ParseFS(content, "template/*.html"))

type PageContent struct {
	CurrentPlayer    *game.Player
	CurrentBlackCard *db.BlackCard
	Sets             map[string]*[]db.Set
}

type JsonAnswerPayload struct {
	Player_id string // can't make this uint without decoding errors....
	Cards     map[int]string
}
type JsonAnswer struct {
	Kind    string
	Winner  string // can't make this uint without decoding errors....
	Payload []JsonAnswerPayload
}

func endRound(w http.ResponseWriter, r *http.Request) {
	var json_answer JsonAnswer
	if err := json.NewDecoder(r.Body).Decode(&json_answer); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 Internal Server Error"))
		return
	}
	// we now have our winner, let's update the game status
	var data game.RoundData
	var err error
	data.Winner, err = strconv.Atoi(json_answer.Winner)
	if err != nil {
		log.Fatal(err)
	}
	data.Submissions = make(map[int][]int, 0)
	for _, p := range json_answer.Payload {
		p_id, err := strconv.Atoi(p.Player_id)
		if err != nil {
			log.Fatal(err)
		}
		data.Submissions[p_id] = make([]int, 0)
		for key := range p.Cards {
			data.Submissions[p_id] = append(data.Submissions[p_id], key)
		}
	}
	game.EndRound(data)
	w.WriteHeader(http.StatusOK)
	// no redirect: JavaScript code will take care of this
}

func serveSetup(w http.ResponseWriter, r *http.Request) {
	pc := PageContent{}
	pc.CurrentBlackCard = nil
	pc.CurrentPlayer = nil
	pc.Sets = make(map[string]*[]db.Set)
	game_status := game.GetGame()
	if len(game_status.Players) == 0 {
		dbs := game.Load()
		for lang, db := range *dbs {
			pc.Sets[lang] = &db.Sets
		}
	} else {
		pc.CurrentBlackCard = game_status.Black_Card
	}
	templates.ExecuteTemplate(w, "setup.html", pc)
}

func joinGame(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.PostForm.Has("player_name") {
		var role string
		if r.PostForm.Has("sets") {
			// first player is joining
			game.Init(r.FormValue("lang"), r.Form["sets"], r.FormValue("player_name"))
			role = "CZAR"
		} else {
			// another player is joining
			game.Join(r.FormValue("player_name"))
			role = "PLAYER"
		}
		game_status := game.GetGame()
		current_player := &game_status.Players[len(game_status.Players)-1]
		log.Printf("%s %d : %s \n", role, current_player.ID, current_player.Name)
		http.Redirect(w, r, "/play/"+strconv.Itoa(current_player.ID)+"/", http.StatusSeeOther)
		return
	} else {
		http.Error(w, "Wrong POST parameters", http.StatusInternalServerError)
		return
	}
}

func serveGame(w http.ResponseWriter, r *http.Request) {
	player_id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		log.Fatal(err)
	}
	game_status := game.GetGame()
	pc := PageContent{}
	pc.CurrentPlayer = &game_status.Players[player_id]
	pc.CurrentBlackCard = game_status.Black_Card

	templates.ExecuteTemplate(w, "game.html", pc)
}

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	mux := http.NewServeMux()
	// pages
	mux.HandleFunc("GET /{$}", serveSetup)
	mux.HandleFunc("GET /setup/", serveSetup)
	mux.HandleFunc("GET /play/{id}/", serveGame)
	// commands
	mux.HandleFunc("POST /join/", joinGame)
	mux.HandleFunc("POST /endround/", endRound)
	// tools
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	// static data
	mux.Handle("/scripts/", http.FileServerFS(scripts))
	mux.Handle("/styles/", http.FileServerFS(styles))
	http.ListenAndServe(":8080", mux)
}
