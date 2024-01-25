package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"

	"github.com/delcaran/cah/db"
)

type player struct {
	id uint
	name string
	score uint
	white_cards []db.WhiteCard
}

type status struct {
	players []player
}

type page_content struct {
	player_num uint
	game_status *status
	czar bool
	black_card db.BlackCard
}


//go:embed template/*
var content embed.FS
var templates = template.Must(template.ParseFS(content, "template/*.html"))
var upgrader = websocket.Upgrader{} // use default options
var game_status status

func echo(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func websocket_test(w http.ResponseWriter, r *http.Request) {
	templates.ExecuteTemplate(w, "websocket_test.html", "ws://"+r.Host+"/echo")
}


func index_handler(w http.ResponseWriter, r *http.Request) {
    err := templates.ExecuteTemplate(w, "main.html", nil)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func new_game_handler(w http.ResponseWriter, r *http.Request) {
	if len(game_status.players) == 0 {
		database, _ := db.Load("eng")
		err := templates.ExecuteTemplate(w, "new.html", database)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

// only the first czar enters here
func card_selection_handler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		log.Fatalf("ParseForm() err: %v", err)
		return
	}
    selected_sets_str := r.Form["sets"] // array of strings
	var selected_sets []int
	var err error
	for _, s := range selected_sets_str {
		v, err := strconv.Atoi(s)
		if err != nil {
			log.Fatal(err)
		}
		selected_sets = append(selected_sets, v)
	}
	db.SelectCards(selected_sets)
	// card selected, go to play
	first_player := player{id:0, score:0, name:""}
	for i:=0; i<10; i++ {
		first_player.white_cards = append(first_player.white_cards, db.GetWhiteCard())
	}
	game_status.players = append(game_status.players, )
	var pc page_content
	pc.game_status = &game_status
	pc.player_num = 0
	pc.czar = true
	pc.black_card = db.GetBlackCard()
	err = templates.ExecuteTemplate(w, "play.html", pc)
	if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

func main() {
	game_status = status{players: make([]player, 0)}
    http.HandleFunc("/", index_handler)
	http.HandleFunc("/new/", new_game_handler)
	http.HandleFunc("/select_sets/", card_selection_handler)
	http.HandleFunc("/echo", echo)
	http.HandleFunc("/test", websocket_test)
    log.Fatal(http.ListenAndServe(":8080", nil))
}
