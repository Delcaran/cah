package db

import (
	_ "embed"
	"encoding/json"
)

//go:embed eng.json
var db_eng string

//go:embed ita.json
var db_ita string

type Card interface {
	Print() string
	getPick() uint
}

type BlackCard struct {
	Text string `json:"text"`
	Pick uint   `json:"pick"`
}

func (c BlackCard) Print() string {
	return c.Text
}
func (c BlackCard) getPick() uint {
	return c.Pick
}


type WhiteCard struct {
	Text string
}

func (c WhiteCard) Print() string {
	return c.Text
}
func (c WhiteCard) getPick() uint {
	return 0
}

var database Database

type Set struct {
	Name string `json:"name"`
	Official bool `json:"official"`
	BlackIDs []uint `json:"black"`
	WhiteIDs []uint `json:"white"`
}

type Database struct {
	Black []BlackCard `json:"black"`
	White []WhiteCard `json:"white"`
	Sets  []Set `json:"sets"`
}

type SelectedCards struct {
	black []Card
	white []Card
}

func Load(lang string) (*Database, error) {
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

func GetSelectedCards(sets []int) (SelectedCards, error) {
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

	var selection SelectedCards
	for _, id := range wanted_blacks {
		selection.black = append(selection.black, database.Black[id])
	}
	for _, id := range wanted_whites {
		selection.white = append(selection.black, database.White[id])
	}

	return selection, nil
}