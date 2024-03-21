package db

import (
	"embed"
	"encoding/json"
	"io/fs"
	"math/rand"
	"strings"
	"time"
)

//go:embed all:*.json
var databases_fs embed.FS

var databases map[string]*Database
var database *Database
var avail_blacks []*BlackCard
var avail_whites []*string
var used_blacks []*BlackCard
var used_whites []*string
var randomizer *rand.Rand

type BlackCard struct {
	Text string `json:"text"`
	Pick uint   `json:"pick"`
}

type Set struct {
	Name     string `json:"name"`
	Official bool   `json:"official"`
	BlackIDs []uint `json:"black"`
	WhiteIDs []uint `json:"white"`
}

type Database struct {
	Black []BlackCard `json:"black"`
	White []string    `json:"white"`
	Sets  []Set       `json:"sets"`
}

func Load() (*map[string]*Database, error) {
	randomizer = rand.New(rand.NewSource(time.Now().UnixNano()))
	databases = make(map[string]*Database)
	fs.WalkDir(databases_fs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			values, _ := databases_fs.ReadFile(path)
			name := strings.Split(path, ".")[0]
			var tmp Database
			json.Unmarshal([]byte(values), &tmp)
			databases[name] = &tmp
		}
		return nil
	})

	return &databases, nil
}

func SelectdDB(lang string) {
	database = databases[lang]
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
	avail_blacks = make([]*BlackCard, 0)
	avail_whites = make([]*string, 0)
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
	if len(sets) > 1 {
		wanted_blacks = cleanup_card_list(wanted_blacks)
		wanted_whites = cleanup_card_list(wanted_whites)
	}

	for _, id := range wanted_blacks {
		avail_blacks = append(avail_blacks, &database.Black[id])
	}
	for _, id := range wanted_whites {
		avail_whites = append(avail_whites, &database.White[id])
	}
	avail_blacks = shuffle[BlackCard](avail_blacks)
	avail_whites = shuffle[string](avail_whites)
}

func shuffle[C BlackCard | string](deck []*C) []*C {
	dest := make([]*C, len(deck))
	perm := rand.Perm(len(deck))
	for i, v := range perm {
		dest[v] = deck[i]
	}
	return dest
}

func getCard[C BlackCard | string](deck []*C, used_deck []*C) (*C, []*C, []*C) {
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

func GetWhiteCard() *string {
	var extracted *string
	extracted, avail_whites, used_whites = getCard[string](avail_whites, used_whites)
	return extracted
}
