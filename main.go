package main

import (
	_ "embed"

	"github.com/delcaran/cah/db"
)

func main() {
	db.Load("eng")
	db.Load("ita")
}
