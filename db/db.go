package db

import (
	_ "embed"
	"encoding/json"
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
	black []BlackCard
	white []WhiteCard
}

func Load(lang string) (*Database, error){
	var values string
	data := Database{}
	switch lang {
	case "ita":
		values = db_ita
	case "eng":
		values = db_eng
	}
    json.Unmarshal([]byte(values), &data)
	return &data, nil
}
