package game

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
	Cards []*string
}

type Status struct {
	Players    []Player
	Black_Card *db.BlackCard
}

var game_status = Status{Players: make([]Player, 0)}

func LoadDatabase(language string) *db.Database {
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

// Load the cards, setup first player and first black card
func Init(selected_sets_str []string, first_czar string) {
	var selected_sets []int
	for _, s := range selected_sets_str {
		v, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		selected_sets = append(selected_sets, v)
	}
	db.SelectCards(selected_sets)
	first_player := Player{ID: 0, Score: 0, Name: first_czar, Czar: true}
	for i := 0; i < 10; i++ {
		first_player.Cards = append(first_player.Cards, db.GetWhiteCard())
	}
	game_status.Players = append(game_status.Players, first_player)
	game_status.Black_Card = db.GetBlackCard()
}

// a new player enters the game
func Join(playername string) {
	new_player := Player{ID: uint(len(game_status.Players)), Score: 0, Name: playername, Czar: false, Cards: make([]*string, 0)}
	for i := 0; i < 10; i++ {
		new_player.Cards = append(new_player.Cards, db.GetWhiteCard())
	}
	game_status.Players = append(game_status.Players, new_player)
}

func GetGame() *Status {
	return &game_status
}
