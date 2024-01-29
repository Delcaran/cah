package db

import (
	_ "embed"
	"encoding/json"
	"math/rand"
	"time"
)

//go:embed eng.json
var db_eng string

//go:embed ita.json
var db_ita string

var database Database // all my sets
var avail_blacks []*BlackCard
var avail_whites []*WhiteCard
var used_blacks []*BlackCard
var used_whites []*WhiteCard
var randomizer *rand.Rand

type BlackCard struct {
	Text string `json:"text"`
	Pick uint   `json:"pick"`
}

type WhiteCard struct {
	Text string
}

type Set struct {
	Name     string `json:"name"`
	Official bool   `json:"official"`
	BlackIDs []uint `json:"black"`
	WhiteIDs []uint `json:"white"`
}

type Database struct {
	Black []BlackCard `json:"black"`
	White []WhiteCard `json:"white"`
	Sets  []Set       `json:"sets"`
}

func Load(lang string) (*Database, error) {
	randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
	var values string
	switch lang {
	case "ita":
		values = db_ita
	case "eng":
		values = db_eng
	}
	json.Unmarshal([]byte(values), &database)
	return &database, nil
}

func cleanup_card_list(cards []uint) []uint {
	var cleaned []uint
	for _, x := range cards {
		present := false
		for _, y := range cleaned {
			if x == y {
				present = true
				break
			}
		}
		if !present {
			cleaned = append(cleaned, x)
		}
	}
	return cleaned
}

func SelectCards(sets []int) {
	var wanted_blacks []uint
	var wanted_whites []uint
	for index, set := range database.Sets {
		for _, sel := range sets {
			if sel == index {
				wanted_blacks = append(wanted_blacks, set.BlackIDs...)
				wanted_whites = append(wanted_whites, set.WhiteIDs...)
			}
		}
	}
	wanted_blacks = cleanup_card_list(wanted_blacks)
	wanted_whites = cleanup_card_list(wanted_whites)

	for _, id := range wanted_blacks {
		avail_blacks = append(avail_blacks, &database.Black[id])
	}
	for _, id := range wanted_whites {
		avail_whites = append(avail_whites, &database.White[id])
	}
	avail_blacks = shuffle[BlackCard](avail_blacks)
	avail_whites = shuffle[WhiteCard](avail_whites)
}

func shuffle[C BlackCard | WhiteCard](deck []*C) []*C {
	dest := make([]*C, len(deck))
	perm := rand.Perm(len(deck))
	for i, v := range perm {
		dest[v] = deck[i]
	}
	return dest
}

func getCard[C BlackCard | WhiteCard](deck []*C, used_deck []*C) (*C, []*C, []*C) {
	var extracted *C
	card_num := len(deck)
	if card_num <= 0 {
		deck = append(deck, used_deck...)
		used_deck = nil
	}
	index := randomizer.Intn(card_num)
	extracted = deck[index]
	deck = append(deck[:index], deck[index:]...)
	used_deck = append(used_deck, extracted)
	return extracted, deck, used_deck
}

func GetBlackCard() *BlackCard {
	var extracted *BlackCard
	extracted, avail_blacks, used_blacks = getCard[BlackCard](avail_blacks, used_blacks)
	return extracted
}

func GetWhiteCard() *WhiteCard {
	var extracted *WhiteCard
	extracted, avail_whites, used_whites = getCard[WhiteCard](avail_whites, used_whites)
	return extracted
}
