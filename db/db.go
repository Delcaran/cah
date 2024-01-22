package db

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed eng.json
var db_eng string

//go:embed ita.json
var db_ita string

type BlackCard struct {
	Text string `json:"text"`
	Pick uint   `json:"pick"`
}

type WhiteCard struct {
	Text string
}

type set struct {
	Name string `json:"name"`
	Official bool `json:"official"`
	BlackIDs []uint `json:"black"`
	WhiteIDs []uint `json:"white"`
}

type database struct {
	Black []BlackCard `json:"black"`
	White []WhiteCard `json:"white"`
	Sets  []set `json:"sets"`
}

type SelectedCards struct {
	black []BlackCard
	white []WhiteCard
}

func Load(lang string) {
	var values string
	data := database{}
	switch lang {
	case "ita":
		values = db_ita
	case "eng":
		values = db_eng
	}
    json.Unmarshal([]byte(values), &data)
	for set_num, set := range data.Sets {
		fmt.Printf("set %d name: %s\n", set_num, set.Name)
	}
}
