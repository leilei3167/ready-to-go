package main

import (
	"log"
	"net/http"
	"ready-to-go/Test/TDD/Practice/http-server/poker"
)

const dbFileName = "game.db.json"

func main() {
	store, close, err := poker.FileSysStoreFromFile(dbFileName)
	if err != nil {
		log.Fatal(err)
	}
	defer close()

	server := poker.NewPlayerServer(store)

	log.Fatal(http.ListenAndServe(":5000", server))

}
