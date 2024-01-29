package main

import (
	"log"
	"strconv"

	"github.com/delcaran/cah/db"
)

type Player struct {
	ID    uint
	Czar  bool
	Name  string
	Score uint
	Cards []db.WhiteCard
}

type Status struct {
	players []Player
}

type PageContent struct {
	CurrPlayer *Player
	Black_Card db.BlackCard
}

var game_status Status

func loadDatabase(language string) *db.Database {
	allowed_langs := []string{"eng", "ita"}
	allowed := false
	for _, x := range allowed_langs {
		if x == language {
			allowed = true
		}
	}
	if allowed {
		database, _ := db.Load(language)
		return database
	}
	return nil
}

func initGame(selected_sets_str []string, first_czar string) PageContent {
	var selected_sets []int
	for _, s := range selected_sets_str {
		v, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		selected_sets = append(selected_sets, v)
	}
	db.SelectCards(selected_sets)
	// card selected, go to play
	first_player := Player{ID: 0, Score: 0, Name: first_czar, Czar: true}
	for i := 0; i < 10; i++ {
		first_player.Cards = append(first_player.Cards, db.GetWhiteCard())
	}
	game_status.players = append(game_status.players, first_player)
	var pc PageContent
	pc.CurrPlayer = &game_status.players[0]
	pc.Black_Card = db.GetBlackCard()
	return pc
}

func newGame() *Status {
	return &Status{players: make([]Player, 0)}
}
