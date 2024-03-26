package game

import (
	"log"
	"strconv"

	"github.com/delcaran/cah/db"
)

type Player struct {
	ID    int
	Czar  bool
	Name  string
	Score uint
	Cards []*string
}

type Status struct {
	Players    []Player
	Black_Card *db.BlackCard
}

type RoundData struct {
	Winner      int
	Submissions map[int][]int // PlayerID -> cards list
}

var game_status = Status{Players: make([]Player, 0)}

func Load() *map[string]*db.Database {
	dbs, _ := db.Load()
	return dbs
}

// Load the cards, setup first player and first black card
func Init(lang string, selected_sets_str []string, first_czar string) {
	db.SelectdDB(lang)
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
	for _, p := range game_status.Players {
		if p.Name == playername {
			return
		}
	}
	new_player := Player{ID: len(game_status.Players), Score: 0, Name: playername, Czar: false, Cards: make([]*string, 0)}
	for i := 0; i < 10; i++ {
		new_player.Cards = append(new_player.Cards, db.GetWhiteCard())
	}
	game_status.Players = append(game_status.Players, new_player)
}

func GetGame() *Status {
	return &game_status
}

func EndRound(data RoundData) {
	for player_id, player_cards := range data.Submissions {
		game_status.Players[player_id].Czar = false
		for card_id := range player_cards {
			game_status.Players[player_id].Cards[card_id] = db.GetWhiteCard()
		}
	}
	game_status.Players[data.Winner].Score += 1
	game_status.Players[data.Winner].Czar = true
	game_status.Black_Card = db.GetBlackCard()
}
